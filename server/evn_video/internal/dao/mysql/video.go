package mysql

import (
	"context"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/user/attention"
	"dragonsss.cn/evn_common/model/user/collect"
	"dragonsss.cn/evn_common/model/user/favorites"
	"dragonsss.cn/evn_common/model/user/record"
	"dragonsss.cn/evn_common/model/video"
	"dragonsss.cn/evn_common/model/video/barrage"
	"dragonsss.cn/evn_common/model/video/like"
	video2 "dragonsss.cn/evn_grpc/video"
	"dragonsss.com/evn_video/internal/database/gorms"
	"gorm.io/gorm"
)

type VideoDao struct {
	conn *gorms.GormConn
}

func NewVideoDao() *VideoDao {
	return &VideoDao{
		conn: gorms.New(),
	}
}

func (v VideoDao) GetVideoBarrageByID(ctx context.Context, id uint32) (*barrage.BarragesList, error) {
	session := v.conn.Session(ctx)
	var barrageList *barrage.BarragesList
	err := session.
		Where("video_id", id).
		Find(&barrageList).
		Error
	if err != nil {
		return nil, err
	}
	return barrageList, nil
}

func (v VideoDao) GetVideoComments(ctx context.Context, req *video2.GetVideoCommentRequest) (*video.VideosContribution, error) {
	session := v.conn.Session(ctx)
	var videoList *video.VideosContribution
	err := session.
		Where("id", req.VideoID).
		Preload("Likes").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("UserInfo").Order("created_at desc").Limit(int(req.PageInfo.Size)).Offset(int((req.PageInfo.Page - 1) * req.PageInfo.Size))
		}).
		Find(&videoList).
		Error
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

func (v VideoDao) GetUserByID(ctx context.Context, uid uint32) (*user.User, error) {
	session := v.conn.Session(ctx)
	var user1 *user.User
	err := session.
		Where("id", uid).
		Find(&user1).
		Error
	if err != nil {
		return nil, err
	}
	return user1, nil
}

func (v VideoDao) FindVideoByID(ctx context.Context, videoID uint32) (*video.VideosContribution, error) {
	session := v.conn.Session(ctx)
	var video1 *video.VideosContribution
	err := session.
		Where("id", videoID).
		Preload("Likes").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("UserInfo").Order("created_at desc")
		}).
		Preload("Barrage").
		Preload("UserInfo").
		Order("created_at desc").
		Find(&video1).
		Error
	if err != nil {
		return nil, err
	}
	return video1, nil
}

func (v VideoDao) WatchVideo(ctx context.Context, videoID uint32) (bool, error) {
	var count int64
	err := v.conn.Session(ctx).
		Model(&video.VideosContribution{}).
		Where("id", videoID).
		Updates(map[string]interface{}{"heat": gorm.Expr("Heat  + ?", 1)}).
		Count(&count).Error //数据库查询
	if err == gorm.ErrRecordNotFound {
		//未查询到对应的信息
		return false, nil
	}
	return count > 0, nil
}

func (v VideoDao) IsAttention(ctx context.Context, uid uint32, id uint) (bool, error) {
	session := v.conn.Session(ctx)
	var at *attention.Attention
	err := session.
		Where(attention.Attention{Uid: uint(uid), AttentionID: id}).
		Find(&at).
		Error
	if err != nil {
		return false, err
	}
	return at.ID > 0, nil
}

func (v VideoDao) IsLike(ctx context.Context, uid uint32, id uint) (bool, error) {
	session := v.conn.Session(ctx)
	var lk *like.Likes
	err := session.
		Where(like.Likes{Uid: uint(uid), VideoID: id}).
		Find(&lk).
		Error
	if err != nil {
		return false, err
	}
	return lk.ID > 0, nil
}

func (v VideoDao) GetFavoritesList(ctx context.Context, uid uint32) (*favorites.FavoriteList, error) {
	session := v.conn.Session(ctx)
	var fl *favorites.FavoriteList
	err := session.
		Where("uid", uid).
		Preload("UserInfo").
		Preload("CollectList").
		Order("created_at desc").
		Find(&fl).
		Error
	if err != nil {
		return nil, err
	}
	return fl, nil
}

func (v VideoDao) FindIsCollectByFavorites(ctx context.Context, id uint32, flIDs []uint) (bool, error) {
	//没创建收藏夹情况直接false
	if len(flIDs) == 0 {
		return false, nil
	}
	var cl *collect.CollectsList
	session := v.conn.Session(ctx)
	err := session.
		Where("video_id", id).
		Where("favorites_id", flIDs).
		Find(&cl).
		Error
	if err != nil {
		return false, err
	}
	if len(*cl) == 0 {
		return false, nil
	}
	return true, nil
}

func (v VideoDao) AddVideoRecord(ctx context.Context, uid uint32, videoID uint32) error {
	session := v.conn.Session(ctx)
	var rd *record.Record
	err := session.
		Where(record.Record{Uid: uint(uid), Type: "video", ToId: uint(videoID)}).
		Find(&rd).
		Error
	if err != nil {
		return err
	}
	if rd.ID <= 0 {
		//创建记录
		rd.Uid = uint(uid)
		rd.Type = "video"
		rd.ToId = uint(videoID)
		return session.Create(rd).Error
	} else {
		//更新记录
		return session.Where("id", rd.ID).Updates(rd).Error
	}
}

func (v VideoDao) GetRecommendList(ctx context.Context) (*video.VideosContributionList, error) {
	session := v.conn.Session(ctx)
	var videoList *video.VideosContributionList
	err := session.
		Preload("Likes").
		Preload("Comments").
		Preload("Barrage").
		Preload("UserInfo").
		Order("created_at desc").
		Limit(7).
		Find(&videoList).
		Error
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

func (v VideoDao) CreateVideoBarrage(ctx context.Context, bg *barrage.Barrage) (bool, error) {
	session := v.conn.Session(ctx)
	err := session.
		Create(&bg).
		Error
	if err != nil {
		return false, err
	}
	return true, nil
}
