package chatUser

import (
	"context"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.com/evn_ws/internal/dao/mysql"
	consts "dragonsss.com/evn_ws/utils"
	"dragonsss.com/evn_ws/utils/chat"
	"dragonsss.com/evn_ws/utils/receive/socket"
	"dragonsss.com/evn_ws/utils/response"
	socket2 "dragonsss.com/evn_ws/utils/response/socket"
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"strconv"
)

type Engine struct {
	//视频房间
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
	Tid      uint
	Socket   *websocket.Conn
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
			//清空未读消息
			wsRepo := mysql.NewWsDao()
			ctx := context.Background()
			err := wsRepo.UnreadEmpty(ctx, registerMsg.UserInfo.ID, registerMsg.Tid)
			//添加在线记录
			if _, ok := chat.Severe.UserMapChannel[registerMsg.UserInfo.ID]; ok {
				//聊天对象在线
				chat.Severe.UserMapChannel[registerMsg.UserInfo.ID].ChatList[registerMsg.Tid] = registerMsg.Socket
			}
			if err != nil {
				zap.L().Error("uid " + strconv.Itoa(int(registerMsg.UserInfo.ID)) + " tid " + strconv.Itoa(int(registerMsg.Tid)) + " 清空未读消息数量失败")
			}

		case cancellationMsg := <-e.Cancellation:
			//删除成员
			delete(e.UserMapChannel, cancellationMsg.UserInfo.ID)
			//删除在线记录
			if _, ok := chat.Severe.UserMapChannel[cancellationMsg.UserInfo.ID]; ok {
				//聊天对象在线
				delete(chat.Severe.UserMapChannel[cancellationMsg.UserInfo.ID].ChatList, cancellationMsg.Tid)
			}
		}
	}
}

func CreateChatByUserSocket(uid uint, tid uint, conn *websocket.Conn) (err error) {
	//创建UserChannel
	userChannel := new(UserChannel)
	//绑定ws
	userChannel.Socket = conn
	wsRepo := mysql.NewWsDao()
	ctx := context.Background()
	user1, err := wsRepo.FindUserById(ctx, int64(uid))
	userChannel.UserInfo = user1
	userChannel.Tid = tid
	userChannel.MsgList = make(chan ChanInfo, 10)

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
		case "sendChatMsgText":
			sendChatMsgText(lre, lre.UserInfo.ID, lre.Tid, info)
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
