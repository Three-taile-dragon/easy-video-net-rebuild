package midd

import (
	"context"
	"dragonsss.cn/evn_api/api/user/rpc"
	"dragonsss.cn/evn_api/config"
	common "dragonsss.cn/evn_common"
	"dragonsss.cn/evn_common/errs"
	"dragonsss.cn/evn_common/model"
	user2 "dragonsss.cn/evn_grpc/user"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"time"
)

func TokenVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1.从Header中获取token
		result := &common.Result{}
		token := c.Request.Header.Get("token")
		//2.调用user服务进行token认证
		ctxo, canel := context.WithTimeout(context.Background(), 2*time.Second)
		defer canel()
		response, err := rpc.UserServiceClient.TokenVerify(ctxo, &user2.TokenRequest{Token: token, Secret: config.C.JC.AccessSecret, IsEncrypt: false})
		//3.处理结果 认证通过，将信息放入gin上下文 失败就返回未登录
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			c.JSON(http.StatusOK, result.Fail(code, msg))
			c.Abort() //防止继续执行
			return
		}
		//成功
		c.Set("currentUid", response.Id)
		c.Set("currentUserName", response.Username)
		c.Next()
	}
}

// TokenVerifyNotNecessary 非必须携带Token
func TokenVerifyNotNecessary() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1.从Header中获取token
		result := &common.Result{}
		token := c.Request.Header.Get("token")
		//用户未登陆时不验证
		if len(token) == 0 {
			c.Next()
			return
		}
		//2.调用user服务进行token认证
		ctxo, canel := context.WithTimeout(context.Background(), 2*time.Second)
		defer canel()
		response, err := rpc.UserServiceClient.TokenVerify(ctxo, &user2.TokenRequest{Token: token, Secret: config.C.JC.AccessSecret, IsEncrypt: false})
		//3.处理结果 认证通过，将信息放入gin上下文 失败就返回未登录
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			c.JSON(http.StatusOK, result.Fail(code, msg))
			c.Abort() //防止继续执行
			return
		}
		//成功
		c.Set("currentUid", response.Id)
		c.Set("currentUserName", response.Username)
		c.Next()
	}
}

func ParameterTokenVerify() gin.HandlerFunc {
	type qu struct {
		Token string `json:"token"`
	}
	return func(c *gin.Context) {
		//1.从Parameter中获取token
		result := &common.Result{}
		req := new(qu)
		if err := c.ShouldBindBodyWith(req, binding.JSON); err != nil {
			code, msg := errs.ParseGrpcError(model.SystemError)
			c.JSON(http.StatusOK, result.Fail(code, msg))
			c.Abort() //防止继续执行
			return
		}
		token := req.Token
		//2.调用user服务进行token认证
		ctxo, canel := context.WithTimeout(context.Background(), 2*time.Second)
		defer canel()
		response, err := rpc.UserServiceClient.TokenVerify(ctxo, &user2.TokenRequest{Token: token, Secret: config.C.JC.AccessSecret, IsEncrypt: false})
		//3.处理结果 认证通过，将信息放入gin上下文 失败就返回未登录
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			c.JSON(http.StatusOK, result.Fail(code, msg))
			c.Abort() //防止继续执行
			return
		}
		//成功
		c.Set("currentUid", response.Id)
		c.Set("currentUserName", response.Username)
		c.Next()
	}
}
