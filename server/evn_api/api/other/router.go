package project

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
}
