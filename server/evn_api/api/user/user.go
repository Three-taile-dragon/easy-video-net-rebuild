package user

import (
	"context"
	"dragonsss.cn/evn_api/pkg/model/user"
	common "dragonsss.cn/evn_common"
	"dragonsss.cn/evn_common/errs"
	"dragonsss.cn/evn_grpc/user/login"
	"dragonsss.cn/evn_user/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type HandleUser struct {
}

func New() *HandleUser {
	return &HandleUser{}
}

//返回验证码 调用grpc

func (*HandleUser) getCaptcha(ctx *gin.Context) {
	result := &common.Result{}
	//获取传入的手机号
	mobile := ctx.PostForm("mobile")
	//对grpc进行两秒超时处理
	c, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	//通过grpc调用 验证码生成函数
	rsp, err := LoginServiceClient.GetCaptcha(c, &login.CaptchaRequest{UserMobile: mobile})
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(rsp.Code))
}

func (h *HandleUser) login(c *gin.Context) {
	result := common.Result{}
	//绑定参数
	var req user.LoginReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//检验
	if err := req.VerifyAccount(); err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, err.Error()))
		return
	}
	//对grpc 进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &login.LoginRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		zap.L().Error("登陆模块结构体赋值出错")
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "系统内部错误"))
		return
	}
	//调用grpc
	loginRsp, err := LoginServiceClient.Login(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//返回结果
	rsp := &user.LoginRsp{}
	err = copier.Copy(rsp, loginRsp)
	if err != nil {
		zap.L().Error("登陆模块返回赋值错误", zap.Error(err))
		c.JSON(http.StatusOK, result.Fail(errs.ParseGrpcError(errs.GrpcError(model.SystemError))))
	}
	c.JSON(http.StatusOK, result.Success(rsp))
}

func (h *HandleUser) register(c *gin.Context) {
	result := common.Result{}
	//绑定参数
	var req user.RegisterReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//参数校验
	if err := req.Verify(); err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, err.Error()))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &login.RegisterRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		zap.L().Error("注册模块结构体赋值出错")
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "系统内部错误"))
		return
	}

	_, err = LoginServiceClient.Register(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.Success(""))
}

func (h *HandleUser) refreshToken(c *gin.Context) {
	result := &common.Result{}
	//获取传入的手机号
	refreshToken := c.PostForm("refreshToken")
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	//通过grpc调用 验证码生成函数
	rrsp, err := LoginServiceClient.RefreshToken(ctx, &login.RefreshTokenRequest{RefreshToken: refreshToken})
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//返回结果
	rsp := &user.TokenList{}
	err = copier.Copy(rsp, rrsp)
	if err != nil {
		zap.L().Error("Token刷新模块返回赋值错误", zap.Error(err))
		c.JSON(http.StatusOK, result.Fail(errs.ParseGrpcError(errs.GrpcError(model.SystemError))))
	}
	c.JSON(http.StatusOK, result.Success(rrsp))
}
