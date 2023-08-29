package article

import (
	"dragonsss.cn/evn_api/api/article/rpc"
	"dragonsss.cn/evn_api/api/cors"
	"dragonsss.cn/evn_api/api/midd"
	"dragonsss.cn/evn_api/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

type RouterArticle struct {
}

func init() {
	log.Println("init article router")
	zap.L().Info("init article router")
	ru := &RouterArticle{}
	router.Register(ru)
}

func (*RouterArticle) Router(r *gin.Engine) {
	//初始化grpc的客户端连接
	rpc.InitRpcArticleClient()
	a := New()
	//路由组
	//不需要登入
	contributionRouterNoVerification := r.Group("/api/contribution")
	contributionRouterNoVerification.Use(cors.Cors())
	{
		contributionRouterNoVerification.POST("/getArticleContributionList", a.getArticleContributionList)
		contributionRouterNoVerification.POST("/getArticleContributionListByUser", a.getArticleContributionListByUser)
		contributionRouterNoVerification.POST("/getArticleComment", a.getArticleComment)
		contributionRouterNoVerification.POST("/getArticleClassificationList", a.getArticleClassificationList)
		contributionRouterNoVerification.POST("/getArticleTotalInfo", a.getArticleTotalInfo)
	}
	//非必须登入
	contributionRouterNotNecessary := r.Group("/api/contribution")
	contributionRouterNotNecessary.Use(cors.Cors())
	contributionRouterNotNecessary.Use(midd.TokenVerifyNotNecessary())
	{
		contributionRouterNotNecessary.POST("/getArticleContributionByID", a.getArticleContributionByID)
	}
	//需要登入
	contributionRouter := r.Group("/api/contribution")
	contributionRouter.Use(cors.Cors())
	contributionRouter.Use(midd.TokenVerify())
	{
		contributionRouter.POST("/createArticleContribution", a.createArticleContribution)
		contributionRouter.POST("/updateArticleContribution", a.updateArticleContribution)
		contributionRouter.POST("/deleteArticleByID", a.deleteArticleByID)
		contributionRouter.POST("/articlePostComment", a.articlePostComment)
		contributionRouter.POST("/getArticleManagementList", a.getArticleManagementList)
	}

}
