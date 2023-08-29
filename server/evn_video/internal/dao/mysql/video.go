package mysql

import (
	"context"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/user/attention"
	"dragonsss.cn/evn_common/model/user/collect"
	"dragonsss.cn/evn_common/model/user/favorites"
	"dragonsss.cn/evn_common/model/user/notice"
	"dragonsss.cn/evn_common/model/user/record"
	"dragonsss.cn/evn_common/model/video"
	"dragonsss.cn/evn_common/model/video/barrage"
	"dragonsss.cn/evn_common/model/video/comments"
	"dragonsss.cn/evn_common/model/video/like"
	video2 "dragonsss.cn/evn_grpc/video"
	"dragonsss.com/evn_video/internal/database"
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

func (v VideoDao) CreateVideo(ctx context.Context, videoContribution *video.VideosContribution) (bool, error) {
	session := v.conn.Session(ctx)
	err := session.
		Create(&videoContribution).
		Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (v VideoDao) UpdateVideo(ctx context.Context, videoContribution *video.VideosContribution) (bool, error) {
	session := v.conn.Session(ctx)
	err := session.
		Where("uid", videoContribution.Uid).
		Where("title", videoContribution.Title).
		Where("created_at", videoContribution.CreatedAt).
		Model(&video.VideosContribution{}).
		Updates(&videoContribution).
		Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (v VideoDao) DeleteVideoByID(ctx context.Context, req *video2.CommonIDAndUIDRequest) (bool, error) {
	session := v.conn.Session(ctx)
	var video1 *video.VideosContribution
	err := session.
		Where("id", req.ID).
		Find(&video1).
		Error
	if err != nil {
		return false, err
	}
	if video1.Uid != uint(req.UID) {
		return false, nil
	}
	err1 := session.
		Where("id", req.ID).
		Delete(&video1).
		Error
	if err1 != nil {
		return false, err1
	}
	return true, nil
}

func (v VideoDao) GetCommentFirstIDByID(ctx context.Context, contentID uint32) (*comments.Comment, error) {
	session := v.conn.Session(ctx)
	var com *comments.Comment
	err := session.
		Where("id", contentID).
		Find(&com).
		Error
	if err != nil {
		return nil, err
	}
	//循环获取最顶层的评论ID
	for com.CommentID != 0 {
		c := context.Background()
		com, err = v.GetCommentFirstIDByID(c, uint32(com.CommentID))
		if err != nil {
			return nil, err
		}
	}
	return com, nil
}

func (v VideoDao) GetCommentUserIDByID(ctx context.Context, contentID uint32) (*comments.Comment, error) {
	session := v.conn.Session(ctx)
	var com *comments.Comment
	err := session.
		Where("id", contentID).
		Find(&com).
		Error
	if err != nil {
		return nil, err
	}
	return com, nil
}

func (v VideoDao) CreateComment(conn database.DbConn, ctx context.Context, comment *comments.Comment) error {
	v.conn = conn.(*gorms.GormConn) //使用事务操作
	return v.conn.Tx(ctx).Create(comment).Error
}

func (v VideoDao) GetVideoManagementList(ctx context.Context, req *video2.GetVideoManagementListRequest) (*video.VideosContributionList, error) {
	session := v.conn.Session(ctx)
	var video1 *video.VideosContributionList
	err := session.
		Where("uid", req.Uid).
		Preload("Likes").
		Preload("Comments").
		Preload("Barrage").
		Limit(int(req.PageInfo.Size)).
		Offset(int((req.PageInfo.Page - 1) * req.PageInfo.Size)).
		Order("created_at desc").
		Find(&video1).
		Error
	if err != nil {
		return nil, err
	}
	return video1, nil
}

func (v VideoDao) GetLikeVideo(conn database.DbConn, ctx context.Context, req *video2.CommonIDAndUIDRequest) (*like.Likes, error) {
	v.conn = conn.(*gorms.GormConn) //使用事务操作
	session := v.conn.Tx(ctx)
	var li *like.Likes
	err := session.
		Where("uid", req.UID).
		Where("video_id", req.ID).
		Find(&li).
		Error
	if err != nil {
		return nil, err
	}
	return li, nil
}

func (v VideoDao) DeleteLikeVideo(conn database.DbConn, ctx context.Context, req *video2.CommonIDAndUIDRequest, li *like.Likes) error {
	v.conn = conn.(*gorms.GormConn) //使用事务操作
	session := v.conn.Tx(ctx)
	err := session.
		Where("uid", req.UID).
		Where("video_id", req.ID).
		Delete(&li).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (v VideoDao) AddNotice(ctx context.Context, uid uint, cid uint, tid uint, tp string, c string) error {
	session := v.conn.Session(ctx)
	err := session.
		Create(&notice.Notice{
			Uid:     uid,
			Cid:     cid,
			ToID:    tid,
			Type:    tp,
			Content: c,
			ISRead:  0,
		}).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (v VideoDao) DeleteNotice(ctx context.Context, videoUid uint, uid uint32, videoID uint32, videoLike string) error {
	session := v.conn.Session(ctx)
	err := session.
		Where(&notice.Notice{
			Uid:  videoUid,
			Cid:  uint(uid),
			ToID: uint(videoID),
			Type: videoLike,
		}).
		Delete(&notice.Notice{}).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (v VideoDao) LikeVideo(conn database.DbConn, ctx context.Context, likeVideo *like.Likes) error {
	v.conn = conn.(*gorms.GormConn) //使用事务操作
	return v.conn.Tx(ctx).Create(likeVideo).Error
}
