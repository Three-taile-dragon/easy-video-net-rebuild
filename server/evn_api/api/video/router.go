package project

import (
	"dragonsss.cn/evn_api/api/cors"
	"dragonsss.cn/evn_api/api/midd"
	"dragonsss.cn/evn_api/api/video/rpc"
	"dragonsss.cn/evn_api/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

type RouterVideo struct {
}

func init() {
	log.Println("init video router")
	zap.L().Info("init video router")
	rv := &RouterVideo{}
	router.Register(rv)
}

func (*RouterVideo) Router(r *gin.Engine) {
	//初始化grpc的客户端连接
	rpc.InitRpcVideoClient()
	v := New()
	//不需要登入
	contributionRouterNoVerification := r.Group("/api/contribution")
	//使用token认证中间件
	contributionRouterNoVerification.Use(cors.Cors())
	{
		contributionRouterNoVerification.GET("/video/barrage/v3/", v.getVideoBarrage)
		contributionRouterNoVerification.GET("/getVideoBarrageList", v.getVideoBarrageList)
		contributionRouterNoVerification.POST("/getVideoComment", v.getVideoComment)
	}
	//非必须登入
	contributionRouterNotNecessary := r.Group("/api/contribution")
	contributionRouterNotNecessary.Use(cors.Cors())
	contributionRouterNotNecessary.Use(midd.TokenVerifyNotNecessary())
	{
		contributionRouterNotNecessary.POST("/getVideoContributionByID", v.getVideoContributionByID)
	}
}
