package ws

import (
	"dragonsss.com/evn_ws/api/midd"
	"dragonsss.com/evn_ws/api/ws/rpc"
	"dragonsss.com/evn_ws/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

type RouterProject struct {
}

func init() {
	log.Println("init ws router")
	zap.L().Info("init ws router")
	ru := &RouterProject{}
	router.Register(ru)
}

func (*RouterProject) Router(r *gin.Engine) {
	//初始化grpc的客户端连接
	rpc.InitRpcProjectClient()
	socketRouter := r.Group("ws").Use(midd.VerificationTokenAsSocket())
	{
		w := New()
		socketRouter.GET("/noticeSocket", w.NoticeSocket)
		socketRouter.GET("/chatSocket", w.ChatSocket)
		socketRouter.GET("/chatUserSocket", w.ChatByUserSocket)
		socketRouter.GET("/liveSocket", w.LiveSocket)
		socketRouter.GET("/videoSocket", w.VideoSocket)
	}
}
