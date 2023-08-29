package mysql

import (
	"context"
	"dragonsss.cn/evn_common/model/article"
	"dragonsss.cn/evn_common/model/liveInfo"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/user/attention"
	"dragonsss.cn/evn_common/model/user/chat/chatList"
	"dragonsss.cn/evn_common/model/user/chat/chatMsg"
	"dragonsss.cn/evn_common/model/user/collect"
	"dragonsss.cn/evn_common/model/user/favorites"
	"dragonsss.cn/evn_common/model/user/notice"
	"dragonsss.cn/evn_common/model/user/record"
	"dragonsss.cn/evn_common/model/video"
	user2 "dragonsss.cn/evn_grpc/user"
	"dragonsss.cn/evn_user/internal/database"
	"dragonsss.cn/evn_user/internal/database/gorms"
	"gorm.io/gorm"
	"time"
)

type UserDao struct {
	conn *gorms.GormConn
}

func NewUserDao() *UserDao {
	return &UserDao{
		conn: gorms.New(),
	}
}

func (u *UserDao) SaveUser(conn database.DbConn, ctx context.Context, mem *user.User) error {
	u.conn = conn.(*gorms.GormConn) //使用事务操作
	return u.conn.Tx(ctx).Create(mem).Error
}

func (u *UserDao) IsExistByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := u.conn.Session(ctx).Model(&user.User{}).Where("email=?", email).Count(&count).Error //数据库查询
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return false, nil
	}
	return count > 0, nil
}

func (u *UserDao) IsExistByNameAndEmail(ctx context.Context, name string) (bool, error) {
	var count int64
	err := u.conn.Session(ctx).Model(&user.User{}).Where("email=? or username=?", name, name).Count(&count).Error //数据库查询
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return false, nil
	}
	return count > 0, nil
}

func (u *UserDao) IsExistByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := u.conn.Session(ctx).Model(&user.User{}).Where("username=?", name).Count(&count).Error //数据库查询
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return false, nil
	}
	return count > 0, nil
}

func (u *UserDao) IsExistByMobile(ctx context.Context, mobile string) (bool, error) {
	var count int64
	err := u.conn.Session(ctx).Model(&user.User{}).Where("mobile=?", mobile).Count(&count).Error //数据库查询
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return false, nil
	}
	return count > 0, nil
}
func (u *UserDao) CheckPassword(ctx context.Context, name string) (*user.User, error) {
	var mem *user.User
	err := u.conn.Session(ctx).Where("username=?", name).First(&mem).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return nil, nil
	}
	return mem, err
}

func (u *UserDao) FindUserByName(ctx context.Context, name string) (*user.User, error) {
	var mem *user.User
	err := u.conn.Session(ctx).Where("username=?", name).First(&mem).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return nil, nil
	}
	return mem, err
}

func (u *UserDao) FindUserById(ctx context.Context, id int64) (*user.User, error) {
	var mem *user.User
	err := u.conn.Session(ctx).Where("id=?", id).First(&mem).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return &user.User{}, nil
	}
	return mem, err
}

func (u *UserDao) UpdateLoginTime(conn database.DbConn, ctx context.Context, name string) error {
	u.conn = conn.(*gorms.GormConn) //使用事务操作
	mem, err := u.FindUserByName(ctx, name)
	if err != nil {
		return err
	}
	mem.UpdatedAt = time.Now()
	// 使用 Updates 方法只更新指定字段，不会插入新的数据
	err = u.conn.Session(ctx).Model(&user.User{}).Where("id = ?", mem.ID).Updates(mem).Error
	return err
}

func (u *UserDao) FindUserByEmail(ctx context.Context, email string) (*user.User, error) {
	var mem *user.User
	err := u.conn.Session(ctx).Where("email=?", email).First(&mem).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return nil, nil
	}
	return mem, err
}

func (u *UserDao) UpdateUser(conn database.DbConn, ctx context.Context, mem *user.User) error {
	u.conn = conn.(*gorms.GormConn) //使用事务操作
	oriMem, err := u.FindUserByEmail(ctx, mem.Email)
	if err != nil {
		return err
	}
	oriMem.Password = mem.Password
	// 使用 Updates 方法只更新指定字段，不会插入新的数据
	err = u.conn.Session(ctx).Model(&user.User{}).Where("id = ?", oriMem.ID).Updates(mem).Error
	return err
}

func (u *UserDao) IsAttention(ctx context.Context, uid uint32, id uint32) (bool, error) {
	var at *attention.Attention
	err := u.conn.Session(ctx).Where("uid =? and attention_id = ?", uid, id).First(&at).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return false, nil
	}
	return true, nil
}

func (u *UserDao) GetAttentionNum(ctx context.Context, id uint32) (*int64, error) {
	var at *attention.Attention
	num := new(int64)
	err := u.conn.Session(ctx).Model(at).Where("uid =? ", id).Count(num).Error
	if err != nil {
		//未查询到对应的信息
		return nil, err
	}
	return num, nil
}

func (u *UserDao) GetVermicelliNum(ctx context.Context, id uint32) (*int64, error) {
	var at *attention.Attention
	num := new(int64)
	err := u.conn.Session(ctx).Model(at).Where("attention_id =? ", id).Count(num).Error
	if err != nil {
		//未查询到对应的信息
		return nil, err
	}
	return num, nil
}

func (u *UserDao) GetVideoListBySpace(ctx context.Context, id uint32) (*video.VideosContributionList, error) {
	session := u.conn.Session(ctx)
	var vvc *video.VideosContributionList
	err := session.
		Preload("Likes").
		Preload("Comments").
		Preload("Barrage").
		Where("uid = ?", id).
		Order("created_at desc").
		Find(&vvc).
		Error
	if err != nil {
		return nil, err
	}
	return vvc, nil
}

func (u *UserDao) GetArticleBySpace(ctx context.Context, id uint32) (*article.ArticlesContributionList, error) {
	session := u.conn.Session(ctx)
	var acl *article.ArticlesContributionList
	err := session.
		Preload("Likes").
		Preload("Classification").
		Preload("Comments").
		Where("uid = ?", id).
		Order("created_at desc").
		Find(&acl).
		Error
	if err != nil {
		return nil, err
	}
	return acl, nil
}

// GetAttentionList 获取关注列表
func (u *UserDao) GetAttentionList(ctx context.Context, id uint32) (*attention.AttentionsList, error) {
	session := u.conn.Session(ctx)
	var att *attention.AttentionsList
	err := session.
		Preload("AttentionUserInfo").
		Where("uid = ?", id).
		Find(&att).
		Error
	if err != nil {
		return nil, err
	}
	return att, nil
}

// GetAttentionListByIdArr 获取关注列表 id数组
func (u *UserDao) GetAttentionListByIdArr(ctx context.Context, id uint32) (arr []uint, err error) {
	arr = make([]uint, 0)
	session := u.conn.Session(ctx)
	var att *attention.AttentionsList
	err = session.
		Where("uid = ?", id).
		Find(&att).
		Error
	if err != nil {
		return nil, err
	}
	//自需要数组
	for _, v := range *att {
		arr = append(arr, v.AttentionID)
	}
	return arr, nil
}

// GetVermicelliList 获取粉丝列表
func (u *UserDao) GetVermicelliList(ctx context.Context, id uint32) (*attention.AttentionsList, error) {
	session := u.conn.Session(ctx)
	var att *attention.AttentionsList
	err := session.
		Preload("UserInfo").
		Where("attention_id = ?", id).
		Find(&att).
		Error
	if err != nil {
		return nil, err
	}
	return att, nil
}

func (u *UserDao) IsExistByID(ctx context.Context, id uint32) (bool, error) {
	var count int64
	err := u.conn.Session(ctx).Model(&user.User{}).Where("id=?", id).Count(&count).Error //数据库查询
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return false, nil
	}
	return count > 0, nil
}

// UpdatePureZero 更新数据存在0值
func (u *UserDao) UpdatePureZero(ctx context.Context, id int64, update map[string]interface{}) (bool, error) {
	err := u.conn.Session(ctx).Model(&user.User{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateUserAvatar 更新用户头像
func (u *UserDao) UpdateUserAvatar(ctx context.Context, tmpUser *user.User) (bool, error) {
	err := u.conn.Session(ctx).Model(&user.User{}).Where("id = ?", tmpUser.ID).Updates(tmpUser).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// FindUserLiveInfo 根据ID查找某人的直播信息
func (u *UserDao) FindUserLiveInfo(ctx context.Context, id int64) (*liveInfo.LiveInfo, error) {
	session := u.conn.Session(ctx)
	var lf *liveInfo.LiveInfo
	err := session.
		Where("uid = ?", id).
		Find(&lf).
		Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return lf, nil
}

func (u *UserDao) IsExistUserLiveInfo(ctx context.Context, liveinfo *liveInfo.LiveInfo) (bool, error) {
	var lf *liveInfo.LiveInfo
	num := new(int64)
	err := u.conn.Session(ctx).Model(lf).Where("uid =? ", liveinfo.Uid).Count(num).Error
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return false, nil
	} else if err != nil {
		return false, err
	}
	if *num >= 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (u *UserDao) SaveUserLiveInfo(ctx context.Context, liveinfo *liveInfo.LiveInfo) (bool, error) {
	session := u.conn.Session(ctx)
	err := session.
		Create(&liveinfo).
		Error
	if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func (u *UserDao) UpdateUserLiveInfo(ctx context.Context, liveinfo *liveInfo.LiveInfo) (bool, error) {
	session := u.conn.Session(ctx)
	err := session.
		Where("uid = ?", liveinfo.Uid).
		Updates(&liveinfo).
		Error
	if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func (u *UserDao) Attention(ctx context.Context, aid uint32, uid uint32) (*attention.Attention, error) {
	session := u.conn.Session(ctx)
	var att *attention.Attention
	err := session.
		Where("uid = ? && attention_id = ?", uid, aid).
		Find(&att).
		Error
	//已关注
	if att.ID > 0 {
		//删除已关注的记录
		err = session.
			Where("id = ?", att.ID).
			Delete(&att).
			Error
	} else {
		//未关注
		err = session.
			Create(&attention.Attention{Uid: uint(uid), AttentionID: uint(aid)}).
			Error
	}
	if err == nil {
		return att, nil
	}
	return nil, err
}

func (u *UserDao) SaveFavorites(ctx context.Context, fs *favorites.Favorites) (bool, error) {
	session := u.conn.Session(ctx)
	err := session.
		Create(&fs).
		Error
	if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func (u *UserDao) FindFavoritesByID(ctx context.Context, id uint32) (*favorites.Favorites, error) {
	session := u.conn.Session(ctx)
	var fs *favorites.Favorites
	err := session.
		Preload("CollectList").
		Where("id = ?", id).
		Order("created_at desc").
		Find(&fs).
		Error
	if err != nil {
		return nil, err
	}
	return fs, nil
}

func (u *UserDao) UpdateFavorities(ctx context.Context, fs *favorites.Favorites) (bool, error) {
	err := u.conn.Session(ctx).Model(&favorites.Favorites{}).Where("id = ?", fs.ID).Updates(fs).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UserDao) GetFavoritesList(ctx context.Context, id uint32) (*favorites.FavoriteList, error) {
	session := u.conn.Session(ctx)
	var fl *favorites.FavoriteList
	err := session.
		Preload("UserInfo").
		Preload("CollectList").
		Where("uid = ?", id).
		Order("created_at desc").
		Find(&fl).
		Error
	if err != nil {
		return nil, err
	}
	return fl, nil
}

func (u *UserDao) DeleteFavorites(ctx context.Context, fs *favorites.Favorites) (bool, error) {
	err := u.conn.Session(ctx).Model(&favorites.Favorites{}).Where("id = ?", fs.ID).Delete(fs).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UserDao) DetectCollectByFavoritesID(ctx context.Context, id uint32) (bool, error) {
	err := u.conn.Session(ctx).Model(&collect.Collect{}).Where("favorites_id = ?", id).Delete(&collect.Collect{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UserDao) SaveCollect(ctx context.Context, cl *collect.Collect) (bool, error) {
	session := u.conn.Session(ctx)
	err := session.
		Create(&cl).
		Error
	if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func (u *UserDao) FindVideoExistWhere(ctx context.Context, videoId uint32) (*collect.CollectsList, error) {
	session := u.conn.Session(ctx)
	var cl *collect.CollectsList
	err := session.
		Where("video_id = ?", videoId).
		Find(&cl).
		Error
	if err != nil {
		return nil, err
	}
	return cl, nil
}

func (u *UserDao) GetVideoInfoByFavoriteID(ctx context.Context, favoriteID uint32) (*collect.CollectsList, error) {
	session := u.conn.Session(ctx)
	var cl *collect.CollectsList
	err := session.
		Preload("VideoInfo").
		Where("favorites_id = ?", favoriteID).
		Find(&cl).
		Error
	if err != nil {
		return nil, err
	}
	return cl, nil
}

func (u *UserDao) GetRecordListByUid(ctx context.Context, req *user2.GetRecordListRequest) (*record.RecordsList, error) {
	session := u.conn.Session(ctx)
	var rl *record.RecordsList
	err := session.
		Preload("VideoInfo.UserInfo").
		Preload("ArticleInfo.UserInfo").
		Preload("Userinfo.LiveInfo").
		Where("uid = ?", req.Uid).
		Limit(int(req.Size)).
		Offset(int((req.Page - 1) * req.Size)).
		Order("created_at desc").
		Find(&rl).
		Error
	if err != nil {
		return nil, err
	}
	return rl, nil
}

func (u *UserDao) ClearRecord(ctx context.Context, id uint32) (bool, error) {
	err := u.conn.Session(ctx).Model(&record.Record{}).Where("uid = ?", id).Delete(&record.Record{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UserDao) DeleteRecordByID(ctx context.Context, req *user2.CommonIDAndUIDRequest) (bool, error) {
	err := u.conn.Session(ctx).Model(&record.Record{}).Where("id = ? and uid = ?", req.ID, req.UID).Delete(&record.Record{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UserDao) GetNoticeList(ctx context.Context, req *user2.GetNoticeListRequest, messageType []string) (*notice.NoticesList, error) {
	session := u.conn.Session(ctx)
	var nl *notice.NoticesList
	if len(messageType) > 0 {
		err := session.
			Preload("VideoInfo").
			Preload("ArticleInfo").
			Preload("UserInfo").
			Where("uid", req.Uid).
			Where("type", messageType).
			Limit(int(req.Size)).
			Offset(int((req.Page - 1) * req.Size)).
			Order("created_at desc").
			Find(&nl).
			Error
		if err != nil {
			return nil, err
		}
		return nl, nil
	} else {
		err := session.
			Preload("VideoInfo").
			Preload("ArticleInfo").
			Preload("UserInfo").
			Where("uid = ?", req.Uid).
			Limit(int(req.Size)).
			Offset(int((req.Page - 1) * req.Size)).
			Order("created_at desc").
			Find(&nl).
			Error
		if err != nil {
			return nil, err
		}
		return nl, nil
	}
}

func (u *UserDao) ReadAllNoticeList(ctx context.Context, req *user2.GetNoticeListRequest) (bool, error) {
	err := u.conn.Session(ctx).Model(&notice.Notice{}).Where(notice.Notice{Uid: uint(req.Uid), ISRead: 0}).Updates(&notice.Notice{ISRead: 1}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UserDao) GetChatListByIO(ctx context.Context, req *user2.CommonIDRequest) (*chatList.ChatList, error) {
	session := u.conn.Session(ctx)
	var cl *chatList.ChatList
	err := session.
		Preload("ToUserInfo").
		Where("uid = ?", req.ID).
		Order("updated_at desc").
		Find(&cl).
		Error
	if err != nil {
		return nil, err
	}
	return cl, nil
}

func (u *UserDao) FindMsgList(ctx context.Context, req *user2.CommonIDRequest, v uint) (*chatMsg.MsgList, error) {
	ids := make([]uint, 0)
	ids = append(ids, uint(req.ID), v)
	session := u.conn.Session(ctx)
	var ml *chatMsg.MsgList
	err := session.
		Preload("UInfo").
		Where("uid", ids).
		Order("created_at desc").
		Limit(30).
		Find(&ml).
		Error
	if err != nil {
		return nil, err
	}
	return ml, nil
}

func (u *UserDao) FindHistoryMsg(ctx context.Context, req *user2.GetChatHistoryMsgRequest) (*chatMsg.MsgList, error) {
	ids := make([]uint, 0)
	ids = append(ids, uint(req.Uid), uint(req.Tid))
	session := u.conn.Session(ctx)
	var ml *chatMsg.MsgList
	err := session.
		Preload("UInfo").
		Preload("TInfo").
		Where("uid", ids).
		Where("tid", ids).
		Where("created_at < ?", req.LastTime).
		Order("created_at desc").
		Limit(30).
		Find(&ml).
		Error
	if err != nil {
		return nil, err
	}
	return ml, nil
}

func (u *UserDao) GetLastMessage(ctx context.Context, req *user2.CommonIDAndUIDRequest) (*chatMsg.Msg, error) {
	session := u.conn.Session(ctx)
	var ms *chatMsg.Msg
	err := session.
		Where("uid = ? or  tid  = ? and tid = ? or tid = ? ", req.UID, req.UID, req.ID, req.ID).
		Order("created_at desc").
		Find(&ms).
		Error
	if err != nil {
		return nil, err
	}
	return ms, nil
}

func (u *UserDao) AddChat(ctx context.Context, ci *chatList.ChatsListInfo) (bool, error) {
	//判断是否存在
	session := u.conn.Session(ctx)
	is := &chatList.ChatsListInfo{}
	is.ID = 0
	err := session.
		Where("uid = ? And tid = ?", ci.Uid, ci.Tid).
		Find(&is).
		Error
	if err != nil {
		return false, err
	}
	if is.ID != 0 {
		//更新
		err := session.
			Model(is).
			Updates(map[string]interface{}{"last_at": is.LastAt, "last_message": is.LastMessage}).
			Error
		if err != nil {
			return false, err
		}
	} else {
		err := session.
			Create(&ci).
			Error
		if err == nil {
			return true, nil
		} else {
			return false, err
		}
	}
	return false, nil
}

func (u *UserDao) DeleteChat(ctx context.Context, req *user2.CommonIDAndUIDRequest) (bool, error) {
	err := u.conn.Session(ctx).Model(&chatMsg.MsgList{}).Where("uid = ? and tid = ?", req.UID, req.ID).Delete(&chatMsg.MsgList{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
