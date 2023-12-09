package other

import (
	"dragonsss.cn/evn_api/api/cors"
	"dragonsss.cn/evn_api/api/midd"
	"dragonsss.cn/evn_api/api/other/rpc"
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
	rpc.InitRpcProjectClient()
	h := New()
	//路由组
	group := r.Group("/api/home")
	////使用token认证中间件
	////group.Use(midd.TokenVerify())
	//group.POST("", h.index)
	group.Use(cors.Cors())
	group.POST("/getHomeInfo", h.getHomeInfo)
	liveRouter := r.Group("/api/live")
	liveRouter.Use(cors.Cors())
	liveRouter.Use(midd.TokenVerify())
	{
		l := NewLive()
		liveRouter.POST("/getLiveRoom", l.getLiveRoom)
		liveRouter.POST("/getLiveRoomInfo", l.getLiveRoomInfo)
		liveRouter.POST("/getBeLiveList", l.getBeLiveList)
	}
	discussRouter := r.Group("/api/contribution")
	discussRouter.Use(cors.Cors())
	discussRouter.Use(midd.TokenVerify())
	{
		d := NewDiscuss()
		discussRouter.POST("/getDiscussVideoList", d.getDiscussVideoList)
		discussRouter.POST("/getDiscussArticleList", d.getDiscussArticleList)
		discussRouter.POST("/getDiscussBarrageList", d.getDiscussBarrageList)
	}
	commonality := r.Group("/api/commonality")
	commonality.Use(cors.Cors())
	commonality.Use(midd.TokenVerify())
	c := NewCommonality()
	{
		//上传
		commonality.POST("/upload", c.upload)
		commonality.POST("/UploadSlice", c.uploadSlice)
		commonality.POST("/uploadCheck", c.uploadCheck)
		commonality.POST("/uploadMerge", c.uploadMerge)
		commonality.POST("/uploadingMethod", c.uploadingMethod)
		commonality.POST("/uploadingDir", c.uploadingDir)
		commonality.POST("/getFullPathOfImage", c.getFullPathOfImage)
		commonality.POST("/uploadOss", c.uploadOss)
	}
	//非必须登入
	contributionRouterNotNecessary := r.Group("/api/commonality")
	contributionRouterNotNecessary.Use(cors.Cors())
	contributionRouterNotNecessary.Use(midd.TokenVerifyNotNecessary())
	{
		contributionRouterNotNecessary.POST("/search", c.search)
	}
}
