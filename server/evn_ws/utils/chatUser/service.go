package chatUser

import (
	"context"
	"dragonsss.cn/evn_common/conversion"
	"dragonsss.cn/evn_common/model/user/chat/chatList"
	"dragonsss.cn/evn_common/model/user/chat/chatMsg"
	"dragonsss.com/evn_ws/config"
	"dragonsss.com/evn_ws/internal/dao"
	"dragonsss.com/evn_ws/internal/dao/mysql"
	"dragonsss.com/evn_ws/internal/database"
	consts "dragonsss.com/evn_ws/utils"
	"dragonsss.com/evn_ws/utils/chat"
	socket2 "dragonsss.com/evn_ws/utils/receive/socket"
	"dragonsss.com/evn_ws/utils/response"
	"dragonsss.com/evn_ws/utils/response/socket"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func sendChatMsgText(ler *UserChannel, uid uint, tid uint, info *socket2.Receive) {

	//添加消息记录
	cm := chatMsg.Msg{
		Uid:     uid,
		Tid:     tid,
		Type:    info.Type,
		Message: info.Data,
	}
	wsRepo := mysql.NewWsDao()
	ctx := context.Background()
	err := AddMessage(&cm)
	if err != nil {
		response.ErrorWs(ler.Socket, "发送失败")
		return
	}
	//消息查询
	msgInfo, err := wsRepo.FindChatMsgByID(ctx, cm.ID)
	if err != nil {
		response.ErrorWs(ler.Socket, "发送消息失败")
		return
	}
	photo, _ := conversion.FormattingJsonSrc(msgInfo.UInfo.Photo, config.C.Host.LocalHost, config.C.Host.TencentOssHost)

	//给自己发消息不推送
	if uid == tid {
		return
	}

	if _, ok := chat.Severe.UserMapChannel[tid]; ok {
		//在线情况
		if _, ok := chat.Severe.UserMapChannel[tid].ChatList[uid]; ok {
			//在与自己聊天窗口 (直接进行推送)
			response.SuccessWs(chat.Severe.UserMapChannel[tid].ChatList[uid], consts.ChatSendTextMsg, socket.ChatSendTextMsgStruct{
				ID:        msgInfo.ID,
				Uid:       msgInfo.Uid,
				Username:  msgInfo.UInfo.Username,
				Photo:     photo,
				Tid:       msgInfo.Tid,
				Message:   msgInfo.Message,
				Type:      msgInfo.Type,
				CreatedAt: msgInfo.CreatedAt,
			})
			return
		} else {
			//添加未读记录
			cl, err := wsRepo.UnreadAutocorrection(ctx, tid, uid)
			if err != nil {
				zap.L().Error("uid " + strconv.Itoa(int(uid)) + " tid " + strconv.Itoa(int(tid)) + " 消息记录自增未读消息数量失败")

			}
			ci, err := wsRepo.FindChatListInfoByID(ctx, uid, tid)
			//推送主socket
			response.SuccessWs(chat.Severe.UserMapChannel[tid].Socket, consts.ChatUnreadNotice, socket.ChatUnreadNoticeStruct{
				Uid:         uid,
				Tid:         tid,
				LastMessage: ci.LastMessage,
				LastMessageInfo: socket.ChatSendTextMsgStruct{
					ID:        msgInfo.ID,
					Uid:       msgInfo.Uid,
					Username:  msgInfo.UInfo.Username,
					Photo:     photo,
					Tid:       msgInfo.Tid,
					Message:   msgInfo.Message,
					Type:      msgInfo.Type,
					CreatedAt: msgInfo.CreatedAt,
				},
				Unread: cl.Unread,
			})
		}
	} else {
		//不在线
		_, err := wsRepo.UnreadAutocorrection(ctx, tid, uid)
		if err != nil {
			zap.L().Error("uid " + strconv.Itoa(int(uid)) + " tid " + strconv.Itoa(int(tid)) + " 消息记录自增未读消息数量失败")

		}
	}
}

func AddMessage(m *chatMsg.Msg) error {
	//使用事务
	t := dao.NewTransaction()
	err := t.Action(func(conn database.DbConn) error {
		//添加记录
		wsRepo := mysql.NewWsDao()
		ctx := context.Background()
		err := wsRepo.CreateMsg(ctx, m)
		if err != nil {
			return fmt.Errorf("添加聊天记录失败")
		}
		//聊天列表内添加记录
		uci := &chatList.ChatsListInfo{
			Uid:         m.Uid,
			Tid:         m.Tid,
			LastMessage: m.Message,
			LastAt:      time.Now(),
		}
		err = wsRepo.AddChat(ctx, uci)
		if err != nil {
			return fmt.Errorf("添加聊天列表记录失败")
		}
		tci := &chatList.ChatsListInfo{
			Uid:         m.Tid,
			Tid:         m.Uid,
			LastMessage: m.Message,
			LastAt:      time.Now(),
		}
		err = wsRepo.AddChat(ctx, tci)
		if err != nil {
			return fmt.Errorf("添加聊天列表记录失败")
		}
		return nil
	})
	return err
}
