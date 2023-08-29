package other_service_v1

import (
	"context"
	"dragonsss.cn/evn_common/conversion"
	"dragonsss.cn/evn_common/errs"
	model2 "dragonsss.cn/evn_common/model"
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/rotograph"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/video"
	"dragonsss.cn/evn_grpc/other"
	"dragonsss.cn/evn_other/config"
	"dragonsss.cn/evn_other/internal/dao"
	"dragonsss.cn/evn_other/internal/dao/mysql"
	"dragonsss.cn/evn_other/internal/database/tran"
	"dragonsss.cn/evn_other/internal/repo"
	"dragonsss.cn/evn_other/pkg/model"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// OtherService grpc 登陆服务 实现
type OtherService struct {
	other.UnimplementedOtherServiceServer
	cache       repo.Cache
	transaction tran.Transaction
	menuRepo    repo.OtherRepo
}

func New() *OtherService {
	return &OtherService{
		cache:       dao.Rc,
		transaction: dao.NewTransaction(),
		menuRepo:    mysql.NewOtherDao(),
	}
}

func (o *OtherService) GetHomeInfo(ctx context.Context, req *other.HomeRequest) (*other.HomeResponse, error) {
	c := context.Background()
	//获取主页轮播图
	var rotographList *rotograph.List
	rotographList, err := o.menuRepo.FindRotographInfo(c)
	if err != nil {
		zap.L().Error("evn_other other_service GetHomeInfo FindRotographInfo DB_error", zap.Error(err))
		return nil, errs.GrpcError(model2.DBError)
	}

	//获取主页推荐视频
	var videoList *video.VideosContributionList
	tmp := common.PageInfo{
		Page: int(req.Page),
		Size: int(req.Size),
	}
	videoList, err = o.menuRepo.FindVideoList(c, tmp)

	if err != nil {
		zap.L().Error("evn_other other_service GetHomeInfo FindVideoList DB_error", zap.Error(err))
		return nil, errs.GrpcError(model2.DBError)
	}
	res := &model.GetHomeInfoResponse{}
	res.Response(rotographList, videoList, config.C.Host.LocalHost, config.C.Host.TencentOssHost)

	rotographJSON, err := json.Marshal(res.Rotograph)
	if err != nil {
		zap.L().Error("evn_other other_service GetHomeInfo rotographJSON error", zap.Error(err))
		return nil, errs.GrpcError(model2.JsonError)
	}
	videoJSON, err := json.Marshal(res.VideoList)
	if err != nil {
		zap.L().Error("evn_other other_service GetHomeInfo videoJSON error", zap.Error(err))
		return nil, errs.GrpcError(model2.JsonError)
	}

	response := &other.HomeResponse{
		Rotograph: string(rotographJSON),
		VideoList: string(videoJSON),
	}
	return response, nil
}

func (o *OtherService) GetLiveRoom(ctx context.Context, req *other.CommonIDRequest) (*other.GetLiveRoomResponse, error) {
	//请求直播服务器
	url := config.C.Live.Agreement + "://" + config.C.Live.IP + ":" + strconv.Itoa(config.C.Live.Api) + "/control/get?room="
	url = url + strconv.Itoa(int(req.ID))
	// 创建http get请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	// 解析http请求中body 数据到我们定义的结构体中
	ReqGetRoom := &model.ReqGetRoom{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(ReqGetRoom); err != nil {
		zap.L().Error("evn_other other_service GetLiveRoom Decode error", zap.Error(err))
		return nil, errs.GrpcError(model2.RequestError)
	}
	if ReqGetRoom.Status != 200 {
		return nil, errs.GrpcError(model2.RequestError)
	}
	response := model.GetLiveRoomResponse("rtmp://"+config.C.Live.IP+":"+strconv.Itoa(config.C.Live.Rtmp)+"/live", ReqGetRoom.Data)
	return &other.GetLiveRoomResponse{
		Address: response.(model.GetLiveRoomResponseStruct).Address,
		Key:     response.(model.GetLiveRoomResponseStruct).Key,
	}, nil
}

func (o *OtherService) GetLiveRoomInfo(ctx context.Context, req *other.CommonIDAndUIDRequest) (*other.CommonDataResponse, error) {
	c := context.Background()
	userInfo, err := o.menuRepo.FindLiveInfo(c, req)
	if err != nil {
		zap.L().Error("evn_other other_service GetLiveRoomInfo FindLiveInfo DB_error", zap.Error(err))
		return nil, errs.GrpcError(model2.DBError)
	}
	flv := config.C.Live.Agreement + "://" + config.C.Live.IP + ":" + strconv.Itoa(config.C.Live.Flv) + "/live/" + strconv.Itoa(int(req.ID)) + ".flv"

	if req.UID > 0 {
		//添加历史记录
		err = o.menuRepo.AddLiveRecord(c, req.UID, req.ID)
		if err != nil {
			zap.L().Error("evn_other other_service GetLiveRoomInfo AddLiveRecord error", zap.Error(err))
			return nil, errs.GrpcError(model2.DBError)
		}
	}
	licenseUrl := config.C.Vod.LicenseUrl
	infoResponse := model.GetLiveRoomInfoResponse(userInfo, flv, licenseUrl, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	rspJSON, err := json.Marshal(infoResponse)
	if err != nil {
		zap.L().Error("evn_other other_service GetLiveRoomInfo Marshal error", zap.Error(err))
		return nil, errs.GrpcError(model2.JsonError)
	}
	tmp := &other.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (o *OtherService) GetBeLiveList(ctx context.Context, req *other.CommonIDRequest) (*other.CommonDataResponse, error) {
	c := context.Background()
	//获取开通播放用户id
	url := config.C.Live.Agreement + "://" + config.C.Live.IP + ":" + strconv.Itoa(config.C.Live.Api) + "/stat/livestat"

	// 创建http get请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	// 解析http请求中body 数据到我们定义的结构体中
	liveStat := &model.LivestatRes{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(liveStat); err != nil {
		zap.L().Error("evn_other other_service GetBeLiveList Decode error", zap.Error(err))
		return nil, errs.GrpcError(model2.RequestError)
	}
	if liveStat.Status != 200 {
		return nil, errs.GrpcError(model2.RequestError)
	}
	//获取live中正在直播列表
	keys := make([]uint, 0)
	for _, kv := range liveStat.Data.Publishers {
		ka := strings.Split(kv.Key, "live/")
		uintKey, _ := strconv.ParseUint(ka[1], 10, 19)
		keys = append(keys, uint(uintKey))
	}
	userList := &user.UserList{}
	if len(keys) > 0 {
		userList, err = o.menuRepo.GetBeLiveList(c, keys)
		if err != nil {
			zap.L().Error("evn_other other_service GetBeLiveList GetBeLiveList DB_error", zap.Error(err))
			return nil, errs.GrpcError(model2.DBError)
		}
	}

	listResponse := model.GetBeLiveListResponse(userList, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	rspJSON, err := json.Marshal(listResponse)
	if err != nil {
		zap.L().Error("evn_other other_service GetBeLiveList Marshal error", zap.Error(err))
		return nil, errs.GrpcError(model2.JsonError)
	}
	tmp := &other.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (o *OtherService) GetDiscussVideoList(ctx context.Context, req *other.CommonDiscussRequest) (*other.CommonDataResponse, error) {
	c := context.Background()
	//获取用户发布的视频
	videoList, err := o.menuRepo.FindDiscussVideoCommentList(c, req.Uid)
	if err != nil {
		zap.L().Error("evn_other other_service GetDiscussVideoList FindDiscussVideoCommentList DB_error", zap.Error(err))
		return nil, errs.GrpcError(model2.DBError)
	}
	if len(*videoList) == 0 {
		return &other.CommonDataResponse{}, nil
	}
	videoIDs := make([]uint, 0)
	for _, v := range *videoList {
		videoIDs = append(videoIDs, v.ID)
	}
	//得到视频评论信息
	cml, err := o.menuRepo.GetVideoCommentListByIDs(c, videoIDs, req)
	if err != nil {
		zap.L().Error("evn_other other_service GetDiscussVideoList GetVideoCommentListByIDs DB_error", zap.Error(err))
		return nil, errs.GrpcError(model2.DBError)
	}
	listResponse := model.GetDiscussVideoListResponse(cml, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	rspJSON, err := json.Marshal(listResponse)
	if err != nil {
		zap.L().Error("evn_other other_service GetDiscussVideoList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model2.JsonError)
	}

	response := &other.CommonDataResponse{
		Data: string(rspJSON),
	}
	return response, nil
}

func (o *OtherService) GetDiscussArticleList(ctx context.Context, req *other.CommonDiscussRequest) (*other.CommonDataResponse, error) {
	c := context.Background()
	//获取用户发布的专栏
	articleList, err := o.menuRepo.FindDiscussArticleCommentList(c, req.Uid)
	if err != nil {
		zap.L().Error("evn_other other_service GetDiscussArticleList FindDiscussArticleCommentList DB_error", zap.Error(err))
		return nil, errs.GrpcError(model2.DBError)
	}
	if len(*articleList) == 0 {
		return &other.CommonDataResponse{}, nil
	}
	articleIDs := make([]uint, 0)
	for _, v := range *articleList {
		articleIDs = append(articleIDs, v.ID)
	}
	//得到文章评论信息
	cml, err := o.menuRepo.GetArticleCommentListByIDs(c, articleIDs, req)
	if err != nil {
		zap.L().Error("evn_other other_service GetDiscussArticleList GetArticleCommentListByIDs DB_error", zap.Error(err))
		return nil, errs.GrpcError(model2.DBError)
	}
	listResponse := model.GetDiscussArticleListResponse(cml, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	rspJSON, err := json.Marshal(listResponse)
	if err != nil {
		zap.L().Error("evn_other other_service GetDiscussArticleList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model2.JsonError)
	}

	response := &other.CommonDataResponse{
		Data: string(rspJSON),
	}
	return response, nil
}

func (o *OtherService) GetDiscussBarrageList(ctx context.Context, req *other.CommonDiscussRequest) (*other.CommonDataResponse, error) {
	c := context.Background()
	//获取用户发布的视频
	videoList, err := o.menuRepo.FindDiscussVideoCommentList(c, req.Uid)
	if err != nil {
		zap.L().Error("evn_other other_service GetDiscussBarrageList FindDiscussVideoCommentList DB_error", zap.Error(err))
		return nil, errs.GrpcError(model2.DBError)
	}
	if len(*videoList) == 0 {
		return &other.CommonDataResponse{}, nil
	}
	videoIDs := make([]uint, 0)
	for _, v := range *videoList {
		videoIDs = append(videoIDs, v.ID)
	}
	//得到视频弹幕信息
	cml, err := o.menuRepo.GetVideoBarrageListByIDs(c, videoIDs, req)
	if err != nil {
		zap.L().Error("evn_other other_service GetDiscussBarrageList GetVideoBarrageListByIDs error", zap.Error(err))
		return nil, errs.GrpcError(model2.DBError)
	}
	listResponse := model.GetDiscussBarrageListResponse(cml, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	rspJSON, err := json.Marshal(listResponse)
	if err != nil {
		zap.L().Error("evn_other other_service GetDiscussBarrageList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model2.JsonError)
	}

	response := &other.CommonDataResponse{
		Data: string(rspJSON),
	}
	return response, nil
}

func (o *OtherService) UploadingMethod(ctx context.Context, req *other.UploadingMethodRequest) (*other.UploadingMethodResponse, error) {
	c := context.Background()
	//获取上传类型
	tp, err := o.menuRepo.FindUploadMethod(c, req.Method)
	if err != nil {
		zap.L().Error("evn_other other_service UploadingMethod FindUploadMethod DB_error", zap.Error(err))
		return nil, errs.GrpcError(model2.DBError)
	}
	rsp := model.UploadingMethodResponse(tp.Method)
	return &other.UploadingMethodResponse{Tp: rsp.(model.UploadingMethodResponseStruct).Tp}, nil
}

func (o *OtherService) UploadingDir(ctx context.Context, req *other.UploadingDirRequest) (*other.UploadingDirResponse, error) {
	c := context.Background()
	//获取上传信息
	tp, err := o.menuRepo.FindUploadMethod(c, req.Interface)
	if err != nil {
		zap.L().Error("evn_other other_service UploadingDir FindUploadMethod DB_error", zap.Error(err))
		return nil, errs.GrpcError(model2.DBError)
	}
	rsp := model.UploadingDirResponse(tp.Path, tp.Quality)
	return &other.UploadingDirResponse{Path: rsp.(model.UploadingDirResponseStruct).Path, Quality: float32(rsp.(model.UploadingDirResponseStruct).Quality)}, nil
}

func (o *OtherService) GetFullPathOfImage(ctx context.Context, req *other.GetFullPathOfImageRequest) (*other.CommonDataResponse, error) {
	//获取完整链接
	path, err := conversion.SwitchIngStorageFun(req.Type, req.Path, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_other other_service GetFullPathOfImage SwitchIngStorageFun DB_error", zap.Error(err))
		return nil, errs.GrpcError(model2.DBError)
	}
	return &other.CommonDataResponse{Data: path}, nil
}

func (o *OtherService) Search(ctx context.Context, req *other.SearchRequest) (*other.CommonDataResponse, error) {
	c := context.Background()
	var rspJSON []byte
	switch req.Type {
	case "video":
		//视频搜索
		list, err := o.menuRepo.SearchVideo(c, req.Page, req.Size, req.Keyword)
		if err != nil {
			zap.L().Error("evn_other other_service Search SearchVideo DB_error", zap.Error(err))
			return nil, errs.GrpcError(model2.DBError)
		}
		//if len(*list) == 0 {
		//	return &other.CommonDataResponse{}, errs.GrpcError(model2.DBError)
		//}
		listResponse, _ := model.SearchVideoResponse(list, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
		rspJSON, err = json.Marshal(listResponse)
		if err != nil {
			zap.L().Error("evn_other other_service Search rspJSON error", zap.Error(err))
			return nil, errs.GrpcError(model2.JsonError)
		}

	case "user":
		//搜素用户
		list, err := o.menuRepo.SearchUser(c, req.Page, req.Size, req.Keyword)
		if err != nil {
			zap.L().Error("evn_other other_service Search SearchUser DB_error", zap.Error(err))
			return nil, errs.GrpcError(model2.DBError)
		}
		//if len(*list) == 0 {
		//	return &other.CommonDataResponse{}, nil
		//}
		aids := make([]uint, 0)
		if req.Uid != 0 {
			//用户登入情况下
			al, err := o.menuRepo.GetAttentionList(c, req.Uid)
			if err != nil {
				zap.L().Error("evn_other other_service Search GetAttentionList DB_error", zap.Error(err))
				return nil, errs.GrpcError(model2.DBError)
			}
			for _, v := range *al {
				aids = append(aids, v.AttentionID)
			}
		}
		listResponse, _ := model.SearchUserResponse(list, aids, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
		rspJSON, err = json.Marshal(listResponse)
		if err != nil {
			zap.L().Error("evn_other other_service Search rspJSON error", zap.Error(err))
			return nil, errs.GrpcError(model2.JsonError)
		}

	}

	response := &other.CommonDataResponse{
		Data: string(rspJSON),
	}
	return response, nil
}
