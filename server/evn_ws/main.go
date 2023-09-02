package main

import (
	srv "dragonsss.cn/evn_common"
	_ "dragonsss.com/evn_ws/api"
	"dragonsss.com/evn_ws/config"
	"dragonsss.com/evn_ws/router"
	_ "dragonsss.com/evn_ws/utils/socket"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//r.Use(logs.GinLogger(), logs.GinRecovery(true)) //接收gin框架默认日志
	router.InitRouter(r) //路由初始化
	//grpc服务注册
	gc := router.RegisterGrpc()
	//grpc服务注册到etcd
	router.RegisterEtcdServer()
	stop := func() {
		gc.Stop()
	}
	if config.C.SC.IsHttps {
		srv.RunWithTLS(r, config.C.SC.Name, config.C.SC.Addr, config.C.SC.Cert, config.C.SC.Key, stop)
	} else {
		srv.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
	}
}
