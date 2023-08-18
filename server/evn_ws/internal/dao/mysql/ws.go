package mysql

import (
	"context"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/user/chat/chatList"
	"dragonsss.cn/evn_common/model/user/chat/chatMsg"
	"dragonsss.cn/evn_common/model/user/notice"
	"dragonsss.com/evn_ws/internal/database/gorms"
	"fmt"
	"gorm.io/gorm"
)

type WsDao struct {
	conn *gorms.GormConn
}

func NewWsDao() *WsDao {
	return &WsDao{
		conn: gorms.New(),
	}
}

func (w *WsDao) FindUserById(ctx context.Context, id int64) (*user.User, error) {
	var mem *user.User
	err := w.conn.Session(ctx).Where("id=?", id).First(&mem).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return &user.User{}, nil
	}
	return mem, err
}

func (w *WsDao) GetUnreadNum(ctx context.Context, id uint) (*int64, error) {
	var n *notice.Notice
	num := new(int64)
	err := w.conn.Session(ctx).Model(n).Where("uid", id).Where("is_read", 0).Count(num).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return nil, err
	}
	return num, err
}

func (w *WsDao) GetUnreadNumber(ctx context.Context, id uint) (*int64, error) {
	var chat *chatList.ChatsListInfo
	num := new(int64)
	err := w.conn.Session(ctx).Model(chat).Select("IFNULL(unread,0) as total_unread").Where("uid", id).Where(chatList.ChatsListInfo{Uid: id}).Scan(num).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return nil, err
	}
	return num, err
}

func (w *WsDao) GetListByID(ctx context.Context, uid uint) (*chatList.ChatList, error) {
	var cl *chatList.ChatList
	err := w.conn.Session(ctx).Where("uid", uid).Preload("ToUserInfo").Order("updated_at desc").First(&cl).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return &chatList.ChatList{}, nil
	}
	return cl, err
}

func (w *WsDao) FindList(ctx context.Context, uid uint, tid uint) (*chatMsg.MsgList, error) {
	ids := make([]uint, 0)
	ids = append(ids, uid, tid)
	var cm *chatMsg.MsgList
	err := w.conn.Session(ctx).Where("uid", ids).Where("tid", ids).Preload("UInfo").Preload("TInfo").Order("created_at desc").Limit(30).Find(&cm).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return &chatMsg.MsgList{}, nil
	}
	return cm, err
}

func (w *WsDao) UnreadEmpty(ctx context.Context, uid uint, tid uint) error {
	var cl *chatList.ChatsListInfo
	err := w.conn.Session(ctx).Where(chatList.ChatsListInfo{Uid: uid, Tid: tid}).Find(&cl).Error
	if err != nil {
		return err
	}
	if cl.ID > 0 {
		cl.Unread = 0
		return w.conn.Session(ctx).Save(cl).Error
	}
	return fmt.Errorf("情况失败")
}

func (w *WsDao) CreateMsg(ctx context.Context, m *chatMsg.Msg) error {
	return w.conn.Tx(ctx).Create(m).Error
}

func (w *WsDao) FindChatMsgByID(ctx context.Context, id uint) (*chatMsg.Msg, error) {
	var cm *chatMsg.Msg
	err := w.conn.Session(ctx).Where("id", id).Preload("UInfo").Preload("TInfo").Find(&cm).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return &chatMsg.Msg{}, nil
	}
	return cm, err
}

func (w *WsDao) UnreadAutocorrection(ctx context.Context, tid uint, uid uint) (*chatList.ChatsListInfo, error) {
	var cl *chatList.ChatsListInfo
	err := w.conn.Session(ctx).Where(chatList.ChatsListInfo{Uid: uid, Tid: tid}).Find(&cl).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return &chatList.ChatsListInfo{}, nil
	}
	return cl, err
}

func (w *WsDao) FindChatListInfoByID(ctx context.Context, uid uint, tid uint) (*chatList.ChatsListInfo, error) {
	var cl *chatList.ChatsListInfo
	err := w.conn.Session(ctx).Where(chatList.ChatsListInfo{Uid: uid, Tid: tid}).Find(&cl).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return &chatList.ChatsListInfo{}, nil
	}
	return cl, err
}

func (w *WsDao) AddChat(ctx context.Context, uci *chatList.ChatsListInfo) error {
	//判断是否存在
	is := &chatList.ChatsListInfo{}
	err := w.conn.Session(ctx).Where("uid = ? And tid = ?", uci.Uid, uci.Tid).Find(is).Error
	if err != nil {
		return err
	}
	if is.ID != 0 {
		//存在即更新
		w.conn.Session(ctx).Model(is).Updates(map[string]interface{}{"last_at": uci.LastAt, "last_message": uci.LastMessage})
		return nil
	} else {
		return w.conn.Session(ctx).Create(uci).Error
	}
}

func (w *WsDao) FindUserList(ctx context.Context) (*[]user.User, error) {
	var userList *[]user.User
	err := w.conn.Session(ctx).Select("id").Find(&userList).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return &[]user.User{}, nil
	}
	return userList, err
}
