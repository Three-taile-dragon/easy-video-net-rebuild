package ws_service_v1

import (
	"context"
	"dragonsss.cn/evn_grpc/ws"
	"dragonsss.com/evn_ws/internal/dao"
	"dragonsss.com/evn_ws/internal/dao/mysql"
	"dragonsss.com/evn_ws/internal/database/tran"
	"dragonsss.com/evn_ws/internal/repo"
)

// WsService grpc 登陆服务 实现
type WsService struct {
	ws.UnimplementedWsServiceServer
	cache       repo.Cache
	transaction tran.Transaction
	wsRepo      repo.WsRepo
}

func New() *WsService {
	return &WsService{
		cache:       dao.Rc,
		transaction: dao.NewTransaction(),
		wsRepo:      mysql.NewWsDao(),
	}
}
func (w *WsService) GetChatList(ctx context.Context, req *ws.CommonIDRequest) (*ws.CommonDataResponse, error) {
	return nil, nil
}
