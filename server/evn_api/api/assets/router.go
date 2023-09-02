package assets

import (
	"dragonsss.cn/evn_api/api/cors"
	"dragonsss.cn/evn_api/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

//配置静态文件路由 与本地上传有关

type RouterAssets struct {
}

func init() {
	log.Println("init assets router")
	zap.L().Info("init assets router")
	ru := &RouterAssets{}
	router.Register(ru)
}

func (*RouterAssets) Router(r *gin.Engine) {
	r.Use(cors.Cors())
	PrivateGroup := r.Group("")
	PrivateGroup.Use()
	{
		//静态资源访问
		r.Static("/assets", "/Initial/assets")
	}

}
