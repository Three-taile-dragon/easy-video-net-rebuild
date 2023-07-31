package midd

import (
	"context"
	"dragonsss.cn/evn_api/api/user"
	"dragonsss.cn/evn_api/config"
	common "dragonsss.cn/evn_common"
	"dragonsss.cn/evn_common/errs"
	"dragonsss.cn/evn_grpc/user/login"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func TokenVerify() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		//1.从Header中获取token
		result := &common.Result{}
		token := c.GetHeader("Authorization")
		//2.调用user服务进行token认证
		ctxo, canel := context.WithTimeout(context.Background(), 2*time.Second)
		defer canel()
		response, err := user.LoginServiceClient.TokenVerify(ctxo, &login.TokenRequest{Token: token, Secret: config.C.JC.AccessSecret, IsEncrypt: true})
		//3.处理结果 认证通过，将信息放入gin上下文 失败就返回未登录
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			c.JSON(http.StatusOK, result.Fail(code, msg))
			c.Abort() //防止继续执行
			return
		}
		//成功
		c.Set("memberId", response.Member.Id)
		c.Next()
	}
}
