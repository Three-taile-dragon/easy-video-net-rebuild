package project_service_v1

import (
	"context"
	"dragonsss.cn/evn_common/errs"
	"dragonsss.cn/evn_grpc/project"
	"dragonsss.cn/evn_project/internal/dao"
	"dragonsss.cn/evn_project/internal/dao/mysql"
	"dragonsss.cn/evn_project/internal/data/menu"
	"dragonsss.cn/evn_project/internal/database/tran"
	"dragonsss.cn/evn_project/internal/repo"
	"dragonsss.cn/evn_project/pkg/model"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

// ProjectService grpc 登陆服务 实现
type ProjectService struct {
	project.UnimplementedProjectServiceServer
	cache       repo.Cache
	transaction tran.Transaction
	menuRepo    repo.MenuRepo
}

func New() *ProjectService {
	return &ProjectService{
		cache:       dao.Rc,
		transaction: dao.NewTransaction(),
		menuRepo:    mysql.NewMenuDao(),
	}
}
func (p *ProjectService) Index(ctx context.Context, req *project.IndexRequest) (*project.IndexResponse, error) {
	c := context.Background()
	pms, err := p.menuRepo.FindMenus(c)
	if err != nil {
		zap.L().Error("首页模块menu数据库存入出错", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	childs := menu.CovertChild(pms)
	var mms []*project.MenuMessage
	err = copier.Copy(&mms, childs)
	if err != nil {
		zap.L().Error("首页模块childs结构体赋值错误", zap.Error(err))
		return nil, errs.GrpcError(model.CopyError)
	}
	return &project.IndexResponse{Menus: mms}, nil
}
