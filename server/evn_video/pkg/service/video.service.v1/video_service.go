package video_service_v1

import (
	"context"
	"dragonsss.cn/evn_common/calculate"
	"dragonsss.cn/evn_common/conversion"
	"dragonsss.cn/evn_common/errs"
	"dragonsss.cn/evn_common/model"
	"dragonsss.cn/evn_common/model/common"
	notice2 "dragonsss.cn/evn_common/model/user/notice"
	video2 "dragonsss.cn/evn_common/model/video"
	"dragonsss.cn/evn_common/model/video/barrage"
	"dragonsss.cn/evn_common/model/video/comments"
	"dragonsss.cn/evn_grpc/video"
	"dragonsss.com/evn_video/config"
	"dragonsss.com/evn_video/internal/dao"
	"dragonsss.com/evn_video/internal/dao/mysql"
	"dragonsss.com/evn_video/internal/database"
	"dragonsss.com/evn_video/internal/database/tran"
	"dragonsss.com/evn_video/internal/repo"
	model2 "dragonsss.com/evn_video/pkg/model"
	"dragonsss.com/evn_video/pkg/utils"
	consts "dragonsss.com/evn_ws/utils"
	"dragonsss.com/evn_ws/utils/notice"
	sokcet "dragonsss.com/evn_ws/utils/video"
	"encoding/json"
	"fmt"
	txCommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentyun/vod-go-sdk"
	"go.uber.org/zap"
	"math"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
	if videoInfo.MediaID != "" {
		pSign, _ := utils.PSignCalculate(videoInfo.MediaID)
		rsp.VideoInfo.PSign = pSign
		rsp.VideoInfo.AppID = strconv.FormatInt(config.C.Vod.Appid, 10)
	}
	// 使用腾讯云 云视立方 播放器 需要 LicenseUrl
	rsp.VideoInfo.LicenseUrl = config.C.Vod.LicenseUrl
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

func (vs *VideoService) SendVideoBarrage(ctx context.Context, req *video.SendVideoBarrageRequest) (*video.CommonDataResponse, error) {
	c := context.Background()
	//保存弹幕
	videoID, _ := strconv.ParseUint(req.ID, 0, 19)
	bg := barrage.Barrage{
		Uid:     uint(req.Uid),
		VideoID: uint(videoID),
		Time:    float64(req.Time),
		Author:  req.Author,
		Type:    uint(req.Type),
		Text:    req.Text,
		Color:   uint(req.Color),
	}
	//获取视频评论
	if is, err := vs.videoRepo.CreateVideoBarrage(c, &bg); !is || err != nil {
		zap.L().Error("evn_video video_service SendVideoBarrage CreateVideoBarrage DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	//socket消息通知
	res := sokcet.ChanInfo{
		Type: consts.VideoSocketTypeResponseBarrageNum,
		Data: nil,
	}
	for _, v := range sokcet.Severe.VideoRoom[uint(videoID)] {
		v.MsgList <- res
	}

	rspJSON, err := json.Marshal(req)
	if err != nil {
		zap.L().Error("evn_video video_service SendVideoBarrage rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &video.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (vs *VideoService) CreateVideoContribution(ctx context.Context, req *video.CreateVideoContributionRequest) (*video.CommonDataResponse, error) {
	c := context.Background()
	//发布视频
	videoSrc, _ := json.Marshal(common.Img{
		Src: req.Video,
		Tp:  req.VideoUploadType,
	})
	coverImg, _ := json.Marshal(common.Img{
		Src: req.Cover,
		Tp:  req.CoverUploadType,
	})
	videoContribution := &video2.VideosContribution{
		Uid:           uint(req.Uid),
		Title:         req.Title,
		Cover:         coverImg,
		Reprinted:     conversion.BoolTurnInt8(req.Reprinted),
		Label:         conversion.MapConversionString(req.Label),
		VideoDuration: req.VideoDuration,
		Introduce:     req.Introduce,
		Heat:          0,
	}
	var width, height int

	// 添加视频时长获取
	if videoContribution.VideoDuration != 0 {
		filePath := ""
		// 假设videoPath是你的视频文件路径
		if req.VideoUploadType == "local" {
			//如果是本地上传
			filePath = config.C.UP.LocalConfig.FileUrl + req.Video
		} else {
			filePath = config.C.UP.TencentConfig.TmpFileUrl + req.Video
		}

		cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			// 处理执行命令失败的情况
			zap.L().Error("视频时长获取失败", zap.Error(err))
			// 这里可以根据实际情况返回错误或者设置默认值
		} else {
			// 将输出转换为字符串并去除换行符
			durationStr := strings.TrimSpace(string(output))
			// 将字符串转换为浮点数
			duration, err := strconv.ParseFloat(durationStr, 64)
			if err != nil {
				// 处理转换失败的情况
				zap.L().Error("视频时长转换失败", zap.Error(err))
				// 这里可以根据实际情况返回错误或者设置默认值
			} else {
				// 将浮点数转换为整数（以秒为单位）
				videoContribution.VideoDuration = int64(int(duration))
			}
		}
	}

	if req.VideoUploadType == "local" {
		//如果是本地上传
		filePath := config.C.UP.LocalConfig.FileUrl + req.Video
		var err1 error
		width, height, err1 = calculate.GetVideoResolution(filePath)
		if err1 != nil {
			zap.L().Error("evn_video video_service CreateVideoContribution GetVideoResolution local err", zap.Error(err1))
			return &video.CommonDataResponse{Data: "获取视频分辨率失败"}, errs.GrpcError(model.SystemError)
		}
		// TODO 优化转码部分 更改分辨率辨认方法
		resolutions := []int{1080, 720, 480, 360}
		if height >= 1080 {
			resolutions = resolutions[1:]
			videoContribution.Video = videoSrc
		} else if height >= 720 && height < 1080 {
			resolutions = resolutions[2:]
			videoContribution.Video720p = videoSrc
		} else if height >= 480 && height < 720 {
			resolutions = resolutions[3:]
			videoContribution.Video480p = videoSrc
		} else if height >= 360 && height < 480 {
			resolutions = resolutions[4:]
			videoContribution.Video360p = videoSrc
		} else {
			return &video.CommonDataResponse{Data: "上传视频分辨率过低"}, errs.GrpcError(model.RequestError)
		}
		//进行视频转码
		go func(width, height int, video *video2.VideosContribution) {
			//本地ffmpeg 处理
			inputFile := filePath
			sr := strings.Split(inputFile, ".")
			for _, r := range resolutions {
				// 计算转码后的宽和高需要取整
				w := int(math.Ceil(float64(r) / float64(height) * float64(width)))
				h := int(math.Ceil(float64(r)))
				if h >= height {
					continue
				}
				extraSuffix := fmt.Sprintf("_output_%dp."+sr[1], r)
				dst := sr[0] + extraSuffix
				cmd := exec.Command("ffmpeg",
					"-hwaccel",
					"cuda",
					"-hwaccel_output_format",
					"cuda",
					"-i",
					inputFile,
					"-vf",
					fmt.Sprintf("scale_cuda=%d:%d", w, h),
					"-c:v",
					"h264_nvenc",
					"-preset",
					"medium",
					"-crf",
					"23",
					"-y",
					dst)
				err := cmd.Run()
				if err != nil {
					zap.L().Error(fmt.Sprintf("视频 :%s :转码%d*%d失败 cmd : %s ,err: ", inputFile, w, h, cmd), zap.Error(err))
					continue
				}
				rv := strings.Split(req.Video, ".")
				src, _ := json.Marshal(common.Img{
					Src: rv[0] + extraSuffix,
					Tp:  "local",
				})
				switch r {
				case 1080:
					videoContribution.Video = src
				case 720:
					videoContribution.Video720p = src
				case 480:
					videoContribution.Video480p = src
				case 360:
					videoContribution.Video360p = src
				}
				if is, err := vs.videoRepo.UpdateVideo(c, videoContribution); !is || err != nil {
					zap.L().Error("evn_video video_service CreateVideoContribution CreateVideo DB_error", zap.Error(err))
					//return nil, errs.GrpcError(model.DBError)
				}
				zap.L().Info(fmt.Sprintf("视频 :%s : 转码%d*%d成功", inputFile, w, h))
			}
			//腾讯云 云点播 使用 转自适应码流 通过FileId进行播放

		}(width, height, videoContribution)
	}

	if req.Media != "" {
		videoContribution.MediaID = req.Media
	}
	var vodErr error
	if req.VideoUploadType == "tencentOss" {
		go func(req1 *video.CreateVideoContributionRequest, err1 error) {
			videoFile := filepath.Base(req1.Video)
			coverFile := filepath.Base(req1.Cover)
			Temporary := filepath.ToSlash(config.C.UP.TencentConfig.TmpFileUrl + "assets/tmp")
			client := &vod.VodUploadClient{}
			client.SecretId = config.C.UP.SecretId
			client.SecretKey = config.C.UP.SecretKey
			vodReq := vod.NewVodUploadRequest()
			vodReq.MediaFilePath = txCommon.StringPtr(Temporary + "/" + videoFile)
			vodReq.CoverFilePath = txCommon.StringPtr(Temporary + "/" + coverFile)

			vodRsp, err := client.Upload("ap-guangzhou", vodReq)
			if err != nil {
				err1 = err
				zap.L().Error("evn_video video_service CreateVideoContribution Upload err", zap.Error(err))
				return
			}
			// 只有腾讯云 云点播的添加 FileId 播放时通过FileId 获取
			videoContribution.MediaID = *vodRsp.Response.FileId
			txSrc, _ := json.Marshal(common.Img{
				Src: *vodRsp.Response.MediaUrl,
				Tp:  "tencentOss",
			})
			videoContribution.Video = txSrc
			// 更新数据库
			if is, err := vs.videoRepo.UpdateVideo(c, videoContribution); !is || err != nil {
				zap.L().Error("evn_video video_service CreateVideoContribution UpdateVideo DB_error", zap.Error(err))
				return
			}
		}(req, vodErr)
	}
	if vodErr != nil {
		return &video.CommonDataResponse{Data: "服务器内部错误"}, errs.GrpcError(model.SystemError)
	}

	//保存到数据库
	if is, err := vs.videoRepo.CreateVideo(c, videoContribution); !is || err != nil {
		zap.L().Error("evn_video video_service CreateVideoContribution CreateVideo DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	return &video.CommonDataResponse{Data: "保存成功"}, nil
}

func (vs *VideoService) UpdateVideoContribution(ctx context.Context, req *video.UpdateVideoContributionRequest) (*video.CommonDataResponse, error) {
	c := context.Background()
	//获取视频信息
	videoInfo, err := vs.videoRepo.FindVideoByID(c, req.ID)
	if err != nil {
		zap.L().Error("evn_video video_service UpdateVideoContribution FindVideoByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if videoInfo.Uid != uint(req.Uid) {
		return &video.CommonDataResponse{Data: "非法操作"}, errs.GrpcError(model.ParamsError)
	}
	coverImg, _ := json.Marshal(common.Img{
		Src: req.Cover,
		Tp:  req.CoverUploadType,
	})
	videoInfo.Cover = coverImg
	videoInfo.Title = req.Title
	videoInfo.Label = conversion.MapConversionString(req.Label)
	videoInfo.Reprinted = conversion.BoolTurnInt8(req.Reprinted)
	videoInfo.Introduce = req.Introduce

	if is, err := vs.videoRepo.UpdateVideo(c, videoInfo); !is || err != nil {
		zap.L().Error("evn_video video_service UpdateVideoContribution UpdateVideo DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	return &video.CommonDataResponse{Data: "更新成功"}, nil
}

func (vs *VideoService) DeleteVideoByID(ctx context.Context, req *video.CommonIDAndUIDRequest) (*video.CommonDataResponse, error) {
	c := context.Background()
	if is, err := vs.videoRepo.DeleteVideoByID(c, req); !is || err != nil {
		zap.L().Error("evn_video video_service DeleteVideoByID DeleteVideoByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &video.CommonDataResponse{Data: "删除成功"}, nil
}

func (vs *VideoService) VideoPostComment(ctx context.Context, req *video.VideoPostCommentRequest) (*video.CommonDataResponse, error) {
	c := context.Background()
	//获取视频信息
	videoInfo, err := vs.videoRepo.FindVideoByID(c, req.VideoID)
	if err != nil {
		zap.L().Error("evn_video video_service VideoPostComment FindVideoByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	CommentFirstID, err := vs.videoRepo.GetCommentFirstIDByID(c, req.ContentID)
	if err != nil {
		zap.L().Error("evn_video video_service VideoPostComment GetCommentFirstIDByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	CommentUserID, err := vs.videoRepo.GetCommentUserIDByID(c, req.ContentID)
	if err != nil {
		zap.L().Error("evn_video video_service VideoPostComment GetCommentUserIDByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	comment := comments.Comment{
		Uid:            uint(req.Uid),
		VideoID:        uint(req.VideoID),
		Context:        req.Content,
		CommentID:      uint(req.ContentID),
		CommentUserID:  CommentUserID.Uid,
		CommentFirstID: CommentFirstID.ID,
	}
	//将存入部分使用事务包裹 使得可以回滚数据库操作
	err = vs.transaction.Action(func(conn database.DbConn) error {
		err = vs.videoRepo.CreateComment(conn, c, &comment)
		if err != nil {
			zap.L().Error("evn_video video_service VideoPostComment CreateComment Tx_DB_error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		//消息通知
		if videoInfo.Uid == comment.Uid {
			return nil
		}
		//添加消息通知
		err = vs.videoRepo.AddNotice(c, videoInfo.Uid, comment.Uid, videoInfo.ID, notice2.VideoComment, comment.Context)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &video.CommonDataResponse{Data: "发布失败"}, err
	}
	//socket推送(在线的情况下)
	if _, ok := notice.Severe.UserMapChannel[videoInfo.UserInfo.ID]; ok {
		userChannel := notice.Severe.UserMapChannel[videoInfo.UserInfo.ID]
		userChannel.NoticeMessage(notice2.VideoComment)
	}
	return &video.CommonDataResponse{Data: "发布成功"}, nil
}

func (vs *VideoService) GetVideoManagementList(ctx context.Context, req *video.GetVideoManagementListRequest) (*video.CommonDataResponse, error) {
	c := context.Background()
	//获取个人发布视频信息
	list, err := vs.videoRepo.GetVideoManagementList(c, req)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoManagementList FindVideoByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	rsp, err := model2.GetVideoManagementListResponse(list, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoManagementList GetVideoManagementListResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}

	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_video video_service GetVideoManagementList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &video.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (vs *VideoService) LikeVideo(ctx context.Context, req *video.CommonIDAndUIDRequest) (*video.CommonDataResponse, error) {
	c := context.Background()
	//获取视频信息
	videoInfo, err := vs.videoRepo.FindVideoByID(c, req.ID)
	if err != nil {
		zap.L().Error("evn_video video_service LikeVideo FindVideoByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	//将存入部分使用事务包裹 使得可以回滚数据库操作
	err = vs.transaction.Action(func(conn database.DbConn) error {
		likeVideo, err2 := vs.videoRepo.GetLikeVideo(conn, c, req)
		if err2 != nil {
			zap.L().Error("evn_video video_service LikeVideo GetLikeVideo Tx_DB_error", zap.Error(err2))

			return err
		}
		if likeVideo.ID > 0 {
			err2 := vs.videoRepo.DeleteLikeVideo(conn, c, req, likeVideo)
			if err2 != nil {
				zap.L().Error("evn_video video_service LikeVideo DeleteLikeVideo Tx_DB_error", zap.Error(err2))
				return err
			}
			//点赞自己作品不进行通知
			if videoInfo.UserInfo.ID == uint(req.UID) {
				return nil
			}
			//删除消息通知
			err = vs.videoRepo.DeleteNotice(c, videoInfo.UserInfo.ID, req.UID, req.ID, notice2.VideoLike)
			if err != nil {
				zap.L().Error("evn_video video_service LikeVideo DeleteNotice Tx_DB_error", zap.Error(err))

				return err
			}
		} else {
			likeVideo.Uid = uint(req.UID)
			likeVideo.VideoID = uint(req.ID)
			err = vs.videoRepo.LikeVideo(conn, c, likeVideo)
			if err != nil {
				zap.L().Error("evn_video video_service LikeVideo LikeVideo Tx_DB_error", zap.Error(err))
				return errs.GrpcError(model.DBError)
			}
			//点赞自己作品不进行通知
			if videoInfo.UserInfo.ID == uint(req.UID) {
				return nil
			}
			//添加消息通知
			err = vs.videoRepo.AddNotice(c, videoInfo.UserInfo.ID, uint(req.UID), uint(req.ID), notice2.VideoLike, "赞了您的作品")
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return &video.CommonDataResponse{Data: "操作失败"}, err
	}
	//socket推送(在线的情况下)
	if _, ok := notice.Severe.UserMapChannel[videoInfo.UserInfo.ID]; ok {
		userChannel := notice.Severe.UserMapChannel[videoInfo.UserInfo.ID]
		userChannel.NoticeMessage(notice2.VideoLike)
	}
	return &video.CommonDataResponse{Data: "操作成功"}, nil
}
