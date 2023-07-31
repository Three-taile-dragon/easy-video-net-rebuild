package mysql

import (
	"context"
	"dragonsss.cn/evn_other/internal/data/common"
	"dragonsss.cn/evn_other/internal/data/rotograph"
	"dragonsss.cn/evn_other/internal/data/video"
	"dragonsss.cn/evn_other/internal/database/gorms"
	"strconv"
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
	var offset int
	if info.Page == 1 {
		info.Size = 11
		offset = (info.Page - 1) * info.Size
	} else {
		offset = (info.Page-2)*info.Size + 11
	}
	session := m.conn.Session(ctx)
	var vvc *video.VideosContributionList
	err := session.Preload("Likes").Preload("Comments").Preload("Barrage").Preload("UserInfo").Where("id > " + strconv.Itoa(3108+offset)).Limit(info.Size).Order("created_at desc").Find(&vvc).Error
	//err := session.Where("id > " + strconv.Itoa(3108+offset)).Limit(info.Size).Error
	return vvc, err
}
