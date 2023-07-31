package user

import (
	"dragonsss.cn/evn_api/config"
	"dragonsss.cn/evn_common/discovery"
	"dragonsss.cn/evn_common/logs"
	"dragonsss.cn/evn_grpc/user/login"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
)

//建立rpc连接

var LoginServiceClient login.LoginServiceClient

func InitRpcUserClient() {
	// grpc连接 etcd
	etcdRegister := discovery.NewResolver(config.C.EC.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	// etcd:/// + 服务名
	conn, err := grpc.Dial("etcd:///user", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Error("grpc连接失败,err: " + err.Error())
		log.Fatalf("grpc连接失败,err: %v \n", err)
	}
	LoginServiceClient = login.NewLoginServiceClient(conn)
}
