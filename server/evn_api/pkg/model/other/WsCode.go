package other

type MyCode int64

const (
	CodeDefault MyCode = 0

	CodeSuccess       MyCode = 200
	CodeInvalidParams MyCode = 201
	CodeNoData        MyCode = 202
	CodeServerBusy    MyCode = 500

	CodeInvalidToken      MyCode = 301
	CodeInvalidAuthFormat MyCode = 302
	CodeNotLogin          MyCode = 303

	CodeTypeError MyCode = 415

	CodePasswordMistake MyCode = 1001
)

var msgFlags = map[MyCode]string{
	CodeDefault:       "请求成功",
	CodeSuccess:       "success",
	CodeInvalidParams: "请求参数错误",
	CodeServerBusy:    "服务繁忙",
	CodeNoData:        "没有数据",

	CodeInvalidToken:      "无效的Token",
	CodeInvalidAuthFormat: "认证格式有误",
	CodeNotLogin:          "未登录",

	CodeTypeError: "类型错误",

	CodePasswordMistake: "密码错误",
}

func (c MyCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}

const (
	/*
		videoSocketTypeError  错误
	*/
	VideoSocketTypeError = "error"
	//VideoSocketTypeNumberOfViewers 返回在线观看人数
	VideoSocketTypeNumberOfViewers = "numberOfViewers"
	//VideoSocketTypeSendBarrage 发送弹幕
	VideoSocketTypeSendBarrage = "sendBarrage"
	//VideoSocketTypeResponseBarrageNum 发送弹幕
	VideoSocketTypeResponseBarrageNum = "responseBarrageNum"
	//NoticeSocketTypeMessage 消息通知
	NoticeSocketTypeMessage = "messageNotice"
	//ChatSendTextMsg 聊天界面发送文本消息
	ChatSendTextMsg = "chatSendTextMsg"
	//ChatUnreadNotice 聊天消息未读通知
	ChatUnreadNotice = "chatUnreadNotice"
	//ChatOnlineUnreadNotice 聊天socket初始化推送未读消息
	ChatOnlineUnreadNotice = "chatOnlineUnreadNotice"
)

const (
	//Error 错误信息
	Error = "error"

	/*
		WebClientBarrageReq  发送弹幕请求数据
		WebClientBarrageRes  发送弹幕响应数据
		WebClientHistoricalBarrageRes 历史弹幕消息
	*/
	WebClientBarrageReq           = "webClientBarrageReq"
	WebClientBarrageRes           = "webClientBarrageRes"
	WebClientHistoricalBarrageRes = "webClientHistoricalBarrageRes"

	/*
		WebClientEnterLiveRoomRes  用户上下线提醒
	*/
	WebClientEnterLiveRoomRes = "webClientEnterLiveRoomRes"
)
