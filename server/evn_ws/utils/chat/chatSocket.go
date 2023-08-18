package chat

import (
	"context"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/user/chat/chatMsg"
	response2 "dragonsss.cn/evn_common/response"
	"dragonsss.com/evn_ws/config"
	"dragonsss.com/evn_ws/internal/dao/mysql"
	consts "dragonsss.com/evn_ws/utils"
	"dragonsss.com/evn_ws/utils/receive/socket"
	"dragonsss.com/evn_ws/utils/response"
	socket2 "dragonsss.com/evn_ws/utils/response/socket"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"strconv"
)

type Engine struct {
	UserMapChannel map[uint]*UserChannel

	Register     chan *UserChannel
	Cancellation chan *UserChannel
}

type ChanInfo struct {
	Type string
	Data interface{}
}

// UserChannel 用户信息
type UserChannel struct {
	UserInfo *user.User
	Socket   *websocket.Conn
	ChatList map[uint]*websocket.Conn
	MsgList  chan ChanInfo
}

var Severe = &Engine{
	UserMapChannel: make(map[uint]*UserChannel, 10),
	Register:       make(chan *UserChannel, 10),
	Cancellation:   make(chan *UserChannel, 10),
}

// Start 启动服务
func (e *Engine) Start() {
	for {
		select {
		//注册事件
		case registerMsg := <-e.Register:
			//添加成员
			e.UserMapChannel[registerMsg.UserInfo.ID] = registerMsg
			//进行未读消息通知
			wsRepo := mysql.NewWsDao()
			ctx := context.Background()
			unreadNum, _ := wsRepo.GetUnreadNumber(ctx, registerMsg.UserInfo.ID)
			if *unreadNum > 0 {
				//存在未读消息 直接推送聊天列表和记录
				list, err := GetChatList(ctx, registerMsg.UserInfo.ID)
				if err != nil {
					fmt.Println("查询错误")
					return
				}
				response.SuccessWs(registerMsg.Socket, consts.ChatOnlineUnreadNotice, list)
			}
		case cancellationMsg := <-e.Cancellation:
			//删除成员
			delete(e.UserMapChannel, cancellationMsg.UserInfo.ID)
		}
	}
}

func CreateChatSocket(uid uint, conn *websocket.Conn) (err error) {
	//创建UserChannel
	userChannel := new(UserChannel)
	//绑定ws
	userChannel.Socket = conn
	wsRepo := mysql.NewWsDao()
	ctx := context.Background()
	user1, err := wsRepo.FindUserById(ctx, int64(uid))
	userChannel.UserInfo = user1
	userChannel.MsgList = make(chan ChanInfo, 10)
	userChannel.ChatList = make(map[uint]*websocket.Conn, 0)

	Severe.Register <- userChannel

	go userChannel.Read()
	go userChannel.Writer()
	return nil

}

// Writer 监听写入数据
func (lre *UserChannel) Writer() {
	for {
		select {
		case msg := <-lre.MsgList:
			response.SuccessWs(lre.Socket, msg.Type, msg.Data)
		}
	}
}

// Read 读取数据
func (lre *UserChannel) Read() {
	//链接断开进行离线
	defer func() {
		Severe.Cancellation <- lre
		err := lre.Socket.Close()
		if err != nil {
			return
		}
	}()
	//监听业务通道
	for {
		//检查通达ping通
		lre.Socket.PongHandler()
		_, text, err := lre.Socket.ReadMessage()
		if err != nil {
			return
		}
		info := new(socket.Receive)
		if err = json.Unmarshal(text, info); err != nil {
			response.ErrorWs(lre.Socket, "消息格式错误")
		}
		switch info.Type {

		}
	}
}

func (lre *UserChannel) NoticeMessage(tp string) {
	//获取未读消息
	wsRepo := mysql.NewWsDao()
	ctx := context.Background()
	num, err := wsRepo.GetUnreadNum(ctx, lre.UserInfo.ID)
	if num == nil || err != nil {
		zap.L().Error("通知id为" + strconv.Itoa(int(lre.UserInfo.ID)) + "用户未读消息失败")
	}
	lre.MsgList <- ChanInfo{
		Type: consts.NoticeSocketTypeMessage,
		Data: socket2.NoticeMessageStruct{
			NoticeType: tp,
			Unread:     num,
		},
	}
}

func GetChatList(ctx context.Context, uid uint) (interface{}, error) {
	//获取消息列表
	wsRepo := mysql.NewWsDao()
	cList, err := wsRepo.GetListByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("查询失败")
	}
	ids := make([]uint, 0)
	for _, v := range *cList {
		ids = append(ids, v.Tid)
	}
	msgList := make(map[uint]*chatMsg.MsgList, 0)
	for _, v := range ids {
		ml, err := wsRepo.FindList(ctx, uid, v)
		if err != nil {
			break
		}
		msgList[v] = ml
	}
	res, err := response2.GetChatListResponse(cList, msgList, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		return nil, err
	}
	return res, nil
}
