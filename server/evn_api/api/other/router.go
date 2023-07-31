package project

import (
	"dragonsss.cn/evn_api/api/midd"
	"dragonsss.cn/evn_api/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

type RouterProject struct {
}

func init() {
	log.Println("init other router")
	zap.L().Info("init other router")
	ru := &RouterProject{}
	router.Register(ru)
}

func (*RouterProject) Router(r *gin.Engine) {
	//初始化grpc的客户端连接
	InitRpcProjectClient()
	h := New()
	//路由组
	//group := r.Group("/home/getHomeInfo")
	////使用token认证中间件
	////group.Use(midd.TokenVerify())
	//group.POST("", h.index)
	r.Use(midd.Cors())
	r.POST("/home/getHomeInfo", h.getHomeInfo)

}
