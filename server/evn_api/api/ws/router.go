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
	log.Println("init project router")
	zap.L().Info("init project router")
	ru := &RouterProject{}
	router.Register(ru)
}

func (*RouterProject) Router(r *gin.Engine) {
	//初始化grpc的客户端连接
	InitRpcProjectClient()
	h := New()
	//路由组
	group := r.Group("/api/index")
	//使用token认证中间件
	group.Use(midd.TokenVerify())
	group.POST("", h.index)

}
