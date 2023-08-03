package other_service_v1

import (
	"context"
	"dragonsss.cn/evn_common/errs"
	model2 "dragonsss.cn/evn_common/model"
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/rotograph"
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
		zap.L().Error("evn_other GetHomeInfo error", zap.Error(err))
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
		zap.L().Error("evn_other GetHomeInfo error", zap.Error(err))
		return nil, errs.GrpcError(model2.DBError)
	}
	res := &model.GetHomeInfoResponse{}
	res.Response(rotographList, videoList, config.C.Host.LocalHost, config.C.Host.TencentOssHost)

	rotographJSON, err := json.Marshal(res.Rotograph)
	if err != nil {
		zap.L().Error("evn_other GetHomeInfo rotographJSON error", zap.Error(err))
		return nil, errs.GrpcError(model2.JsonError)
	}
	videoJSON, err := json.Marshal(res.VideoList)
	if err != nil {
		zap.L().Error("evn_other GetHomeInfo videoJSON error", zap.Error(err))
		return nil, errs.GrpcError(model2.JsonError)
	}

	response := &other.HomeResponse{
		Rotograph: string(rotographJSON),
		VideoList: string(videoJSON),
	}
	return response, nil
}
