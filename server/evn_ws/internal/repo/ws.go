package repo

import (
	"context"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/user/chat/chatList"
	"dragonsss.cn/evn_common/model/user/chat/chatMsg"
)

type WsRepo interface {
	FindUserById(ctx context.Context, id int64) (*user.User, error)
	GetUnreadNum(ctx context.Context, id uint) (*int64, error)
	GetUnreadNumber(ctx context.Context, id uint) (*int64, error)
	GetListByID(ctx context.Context, uid uint) (*chatList.ChatList, error)
	FindList(ctx context.Context, uid uint, tid uint) (*chatMsg.MsgList, error)
	UnreadEmpty(ctx context.Context, uid uint, tid uint) error
	CreateMsg(ctx context.Context, m *chatMsg.Msg) error
	FindChatMsgByID(ctx context.Context, id uint) (*chatMsg.Msg, error)
	UnreadAutocorrection(ctx context.Context, tid uint, uid uint) (*chatList.ChatsListInfo, error)
	FindChatListInfoByID(ctx context.Context, uid uint, tid uint) (*chatList.ChatsListInfo, error)
	AddChat(ctx context.Context, uci *chatList.ChatsListInfo) error
}
