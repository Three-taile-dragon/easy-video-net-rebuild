package video_service_v1

import (
	"context"
	"dragonsss.cn/evn_common/errs"
	"dragonsss.cn/evn_common/model"
	"dragonsss.cn/evn_grpc/video"
	"dragonsss.com/evn_video/config"
	"dragonsss.com/evn_video/internal/dao"
	"dragonsss.com/evn_video/internal/dao/mysql"
	"dragonsss.com/evn_video/internal/database/tran"
	"dragonsss.com/evn_video/internal/repo"
	model2 "dragonsss.com/evn_video/pkg/model"
	consts "dragonsss.com/evn_ws/utils"
	"encoding/json"
	"go.uber.org/zap"
	"strconv"
)

// VideoService grpc 登陆服务 实现
type VideoService struct {
	video.UnimplementedVideoServiceServer
	cache       repo.Cache
	transaction tran.Transaction
	videoRepo   repo.VideoRepo
}

func New() *VideoService {
	return &VideoService{
		cache:       dao.Rc,
		transaction: dao.NewTransaction(),
		videoRepo:   mysql.NewVideoDao(),
	}
}

func (vs *VideoService) GetVideoBarrage(ctx context.Context, req *video.CommonIDRequest) (*video.CommonDataResponse, error) {
	c := context.Background()
	//获取视频弹幕
	list, err := vs.videoRepo.GetVideoBarrageByID(c, req.ID)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoBarrage GetVideoBarrageByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	rsp := model2.GetVideoBarrageResponse(list)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoBarrage GetVideoBarrageResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoBarrage rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &video.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (vs *VideoService) GetVideoBarrageList(ctx context.Context, req *video.CommonIDRequest) (*video.CommonDataResponse, error) {
	c := context.Background()
	//获取视频弹幕
	list, err := vs.videoRepo.GetVideoBarrageByID(c, req.ID)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoBarrageList GetVideoBarrageByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	rsp := model2.GetVideoBarrageListResponse(list)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoBarrageList GetVideoBarrageListResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoBarrageList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &video.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (vs *VideoService) GetVideoComment(ctx context.Context, req *video.GetVideoCommentRequest) (*video.CommonDataResponse, error) {
	c := context.Background()
	//获取视频评论
	list, err := vs.videoRepo.GetVideoComments(c, req)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoComment GetVideoComments DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	rsp := model2.GetVideoContributionCommentsResponse(list, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoComment GetVideoBarrageListResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoComment rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &video.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (vs *VideoService) GetVideoContributionByID(ctx context.Context, req *video.GetVideoContributionByIDRequest) (*video.CommonDataResponse, error) {
	c := context.Background()
	//获取视频信息
	videoInfo, err := vs.videoRepo.FindVideoByID(c, req.VideoID)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoContributionByID FindVideoByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	isAttention := false
	isLike := false
	isCollect := false
	rc := dao.Rc
	if req.Uid != 0 {
		//进行视频播放增加
		if !rc.R().SIsMember(c, consts.VideoWatchByID+strconv.Itoa(int(req.VideoID)), req.Uid).Val() {
			//最近无播放
			rc.R().SAdd(c, consts.VideoWatchByID+strconv.Itoa(int(req.VideoID)), req.Uid)
			if is, err := vs.videoRepo.WatchVideo(c, req.VideoID); !is || err != nil {
				zap.L().Error("evn_video video_service GetVideoContributionByID WatchVideo DB_error", zap.Error(err))
			}
			videoInfo.Heat++
		}
		//获取是否关注
		isAttention, err = vs.videoRepo.IsAttention(c, req.Uid, videoInfo.UserInfo.ID)
		if err != nil {
			zap.L().Error("evn_video video_service GetVideoContributionByID IsAttention DB_error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		//获取是否点赞
		isLike, err = vs.videoRepo.IsLike(c, req.Uid, videoInfo.ID)
		if err != nil {
			zap.L().Error("evn_video video_service GetVideoContributionByID IsLike DB_error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		//判断是否已经收藏
		fl, err := vs.videoRepo.GetFavoritesList(c, req.Uid)
		if err != nil {
			zap.L().Error("evn_video video_service GetVideoContributionByID GetFavoritesList DB_error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		flIDs := make([]uint, 0)
		for _, v := range *fl {
			flIDs = append(flIDs, v.ID)
		}
		//判断是否在收藏夹内
		isCollect, err = vs.videoRepo.FindIsCollectByFavorites(c, req.VideoID, flIDs)
		//添加历史记录
		err = vs.videoRepo.AddVideoRecord(c, req.Uid, req.VideoID)
		if err != nil {
			zap.L().Error("evn_video video_service GetVideoContributionByID AddVideoRecord DB_error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
	}
	//获取推荐列表
	recommendList, err := vs.videoRepo.GetRecommendList(c)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoContributionByID GetRecommendList DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	rsp := model2.GetVideoContributionByIDResponse(videoInfo, recommendList, isAttention, isLike, isCollect, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoContributionByID GetVideoContributionByIDResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoContributionByID rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &video.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}
