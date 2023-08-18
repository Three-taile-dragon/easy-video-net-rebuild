package router

import (
	"dragonsss.cn/evn_common/discovery"
	"dragonsss.cn/evn_common/logs"
	"dragonsss.cn/evn_grpc/ws"
	"dragonsss.com/evn_ws/config"
	ws_service_v1 "dragonsss.com/evn_ws/pkg/service/ws.service.v1"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
)

type Router interface {
	Router(r *gin.Engine)
}

// 路由组
var routers []Router

// InitRouter 路由初始化
func InitRouter(r *gin.Engine) {
	for _, reg := range routers {
		reg.Router(r) //注册路由
	}
}

// Register 添加到路由列表中去
func Register(ro ...Router) {
	routers = append(routers, ro...)
}

type gRPCConfig struct {
	Addr         string
	RegisterFunc func(*grpc.Server)
}

// RegisterGrpc 注册grpc服务
func RegisterGrpc() *grpc.Server {
	c := gRPCConfig{
		Addr: config.C.GC.Addr,
		RegisterFunc: func(g *grpc.Server) {
			ws.RegisterWsServiceServer(g, ws_service_v1.New())
		},
	}
	s := grpc.NewServer()
	c.RegisterFunc(s)
	lis, err := net.Listen("tcp", config.C.GC.Addr)
	if err != nil {
		zap.L().Error("grpc地址监听失败")
		log.Println("grpc地址监听失败")
	}
	go func() {
		err = s.Serve(lis)
		if err != nil {
			zap.L().Error("grpc启动失败,err: " + err.Error())
			log.Println("grpc启动失败", err)
			return
		}
	}()
	return s
}

// 注册etcd服务
func RegisterEtcdServer() {
	etcdRegister := discovery.NewResolver(config.C.EC.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	info := discovery.Server{
		Name:    config.C.GC.Name,
		Addr:    config.C.GC.Addr,
		Version: config.C.GC.Version,
		Weight:  config.C.GC.Weight,
	}
	r := discovery.NewRegister(config.C.EC.Addrs, logs.LG)
	_, err := r.Register(info, 2)
	if err != nil {
		zap.L().Error("etcd服务注册失败,err: " + err.Error())
		log.Fatalf("etcd服务注册失败,err: %v \n", err)
	}
}
