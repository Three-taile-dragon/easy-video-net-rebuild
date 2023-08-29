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
	//需要登入 body参数携带 token
	contributionRouterParameter := r.Group("/api/contribution")
	contributionRouterParameter.Use(cors.Cors())
	contributionRouterParameter.Use(midd.ParameterTokenVerify())
	{
		contributionRouterParameter.POST("/video/barrage/v3/", v.sendVideoBarrage)
	}
	//请求头携带
	contributionRouter := r.Group("/api/contribution")
	contributionRouter.Use(cors.Cors())
	contributionRouter.Use(midd.TokenVerify())
	{
		contributionRouter.POST("/createVideoContribution", v.createVideoContribution)
		contributionRouter.POST("/updateVideoContribution", v.updateVideoContribution)
		contributionRouter.POST("/deleteVideoByID", v.deleteVideoByID)
		contributionRouter.POST("/videoPostComment", v.videoPostComment)
		contributionRouter.POST("/getVideoManagementList", v.getVideoManagementList)
		contributionRouter.POST("/likeVideo", v.likeVideo)
	}
}
