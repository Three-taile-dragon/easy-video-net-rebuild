package main

import (
	_ "dragonsss.cn/evn_api/api"
	"dragonsss.cn/evn_api/config"
	"dragonsss.cn/evn_api/router"
	srv "dragonsss.cn/evn_common"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//r.Use(logs.GinLogger(), logs.GinRecovery(true)) //接收gin框架默认日志
	//注册路由
	router.InitRouter(r)
	if config.C.SC.IsHttps {
		srv.RunWithTLS(r, config.C.SC.Name, config.C.SC.Addr, config.C.SC.Cert, config.C.SC.Key, nil)
	} else {
		srv.Run(r, config.C.SC.Name, config.C.SC.Addr, nil)
	}

}
