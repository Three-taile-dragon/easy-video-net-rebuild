package user

import (
	"dragonsss.cn/evn_api/api/cors"
	"dragonsss.cn/evn_api/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

type RouterUser struct {
}

func init() {
	log.Println("init user router")
	zap.L().Info("init user router")
	ru := &RouterUser{}
	router.Register(ru)
}

func (*RouterUser) Router(r *gin.Engine) {
	//初始化grpc客户端连接
	//使得可以调用逻辑函数
	InitRpcUserClient()
	h := New()
	group := r.Group("/api/user")
	group.Use(cors.Cors())
	group.POST("/getCaptcha", h.getCaptcha)
	group.POST("/login", h.login)
	group.POST("/register", h.register)
	group.POST("/refreshToken", h.refreshToken)
}
