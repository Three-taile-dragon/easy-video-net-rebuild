package router

import "github.com/gin-gonic/gin"

type Router interface {
	Router(r *gin.Engine)
}

// 路由组
var routers []Router

// InitRouter 路由初始化
func InitRouter(r *gin.Engine) {
	for _, reg := range routers {
		reg.Router(r) //注册路由
	}
}

// Register 添加到路由列表中去
func Register(ro ...Router) {
	routers = append(routers, ro...)
}
