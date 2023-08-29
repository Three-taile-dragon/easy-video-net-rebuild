package rpc

import (
	"dragonsss.cn/evn_api/config"
	"dragonsss.cn/evn_common/discovery"
	"dragonsss.cn/evn_common/logs"
	"dragonsss.cn/evn_grpc/other"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
)

var OtherServiceClient other.OtherServiceClient

// InitRpcUserClient 初始化grpc客户段连接
func InitRpcProjectClient() {
	//grpc 连接 etcd
	etcdRegister := discovery.NewResolver(config.C.EC.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	// etcd:/// + grpc服务名
	conn, err := grpc.Dial("etcd:///other", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	OtherServiceClient = other.NewOtherServiceClient(conn)
}
