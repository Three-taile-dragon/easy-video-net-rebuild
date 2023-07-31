package user

import (
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
	r.POST("/api/getCaptcha", h.getCaptcha)
	r.POST("/api/login", h.login)
	r.POST("/api/register", h.register)
	r.POST("/api/refreshToken", h.refreshToken)
}
