package mysql

import (
	"context"
	"dragonsss.cn/evn_common/model/article"
	"dragonsss.cn/evn_common/model/liveInfo"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/user/attention"
	"dragonsss.cn/evn_common/model/video"
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
