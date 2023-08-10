package mysql

import (
	"context"
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/rotograph"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/user/record"
	"dragonsss.cn/evn_common/model/video"
	"dragonsss.cn/evn_grpc/other"
	"dragonsss.cn/evn_other/internal/database/gorms"
)

type OtherDao struct {
	conn *gorms.GormConn
}

func NewOtherDao() *OtherDao {
	return &OtherDao{
		conn: gorms.New(),
	}
}

// FindRotographInfo  获取主页轮播图
func (m *OtherDao) FindRotographInfo(ctx context.Context) (*rotograph.List, error) {
	session := m.conn.Session(ctx)
	var rrg *rotograph.List
	err := session.Find(&rrg).Error
	return rrg, err
}

// FindVideoList 实现视频列表查询
func (m *OtherDao) FindVideoList(ctx context.Context, info common.PageInfo) (*video.VideosContributionList, error) {
	// 首页加载13个铺满后续15个
	var offset int
	if info.Page == 1 {
		info.Size = 13
	} else {
		offset = (info.Page-2)*info.Size + 13
	}

	session := m.conn.Session(ctx)
	var vvc video.VideosContributionList // 修改这里，不再是指针类型
	var totalCount int64
	err := session.
		Model(&video.VideosContribution{}).
		Count(&totalCount).
		Error

	if err != nil {
		return nil, err
	}
	// TODO 查询有问题
	err = session.
		Preload("Likes").
		Preload("Comments").
		Preload("Barrage").
		Preload("UserInfo").
		Where("id > ?", 3108+offset).
		Limit(info.Size).
		Order("created_at desc").
		Find(&vvc).
		Error
	return &vvc, err
}

func (m *OtherDao) FindLiveInfo(ctx context.Context, req *other.CommonIDAndUIDRequest) (*user.User, error) {
	session := m.conn.Session(ctx)
	var us *user.User
	err := session.
		Preload("LiveInfo").
		Where("id", req.ID).
		Find(&us).
		Error
	if err != nil {
		return nil, err
	}
	return us, nil
}

func (m *OtherDao) AddLiveRecord(ctx context.Context, uid uint32, id uint32) error {
	session := m.conn.Session(ctx)
	var rd *record.Record
	err := session.
		Where(record.Record{Uid: uint(uid), Type: "live", ToId: uint(id)}).
		Find(&rd).
		Error
	if err != nil {
		return err
	}
	if rd.ID <= 0 {
		//创建记录
		rd.Uid = uint(uid)
		rd.Type = "live"
		rd.ToId = uint(id)
		return session.Create(&rd).Error
	} else {
		//更新记录
		return session.Where("id", rd.ID).Updates(&rd).Error
	}
}

func (m *OtherDao) GetBeLiveList(ctx context.Context, keys []uint) (*user.UserList, error) {
	session := m.conn.Session(ctx)
	var userList *user.UserList
	err := session.
		Preload("LiveInfo").
		Where("id", keys).
		Find(&userList).
		Error
	if err != nil {
		return nil, err
	}
	return userList, nil
}
