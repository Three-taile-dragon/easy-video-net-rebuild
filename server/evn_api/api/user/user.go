package user

import (
	"context"
	"dragonsss.cn/evn_api/api/user/rpc"
	"dragonsss.cn/evn_api/pkg/model"
	modelUser "dragonsss.cn/evn_api/pkg/model/user"
	common "dragonsss.cn/evn_common"
	"dragonsss.cn/evn_common/errs"
	"dragonsss.cn/evn_grpc/user"
	"encoding/json"
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
	//获取传入的邮箱
	//绑定参数
	var req modelUser.EmailCaptcha
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	c, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	//通过grpc调用 验证码生成函数
	rsp, err := rpc.UserServiceClient.GetCaptcha(c, &user.CaptchaRequest{Email: req.Email})
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (h *HandleUser) login(c *gin.Context) {
	result := common.Result{}
	//绑定参数
	var req modelUser.LoginReq
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
	msg := &user.LoginRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		zap.L().Error("登陆模块结构体赋值出错")
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "系统内部错误"))
		return
	}
	//调用grpc
	loginRsp, err := rpc.UserServiceClient.Login(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//返回结果
	//rsp := &modelUser.LoginRsp{}
	//err = copier.Copy(rsp, loginRsp)
	//if err != nil {
	//	zap.L().Error("登陆模块返回赋值错误", zap.Error(err))
	//	c.JSON(http.StatusOK, result.Fail(errs.ParseGrpcError(errs.GrpcError(model.SystemError))))
	//}
	c.JSON(http.StatusOK, result.Success(loginRsp))
}

func (h *HandleUser) register(c *gin.Context) {
	result := common.Result{}
	//绑定参数
	var req modelUser.RegisterReq
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
	msg := &user.RegisterRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		zap.L().Error("evn_api api user register copy RegisterRequest error")
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "系统内部错误"))
		return
	}

	rsp, err := rpc.UserServiceClient.Register(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.Success(rsp))
}

func (h *HandleUser) refreshToken(c *gin.Context) {
	result := &common.Result{}
	//获取传入的手机号
	refreshToken := c.PostForm("refreshToken")
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	//通过grpc调用 验证码生成函数
	rrsp, err := rpc.UserServiceClient.RefreshToken(ctx, &user.RefreshTokenRequest{RefreshToken: refreshToken})
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//返回结果
	rsp := &modelUser.TokenList{}
	err = copier.Copy(rsp, rrsp)
	if err != nil {
		zap.L().Error("Token刷新模块返回赋值错误", zap.Error(err))
		c.JSON(http.StatusOK, result.Fail(errs.ParseGrpcError(errs.GrpcError(model.SystemError))))
	}
	c.JSON(http.StatusOK, result.Success(rrsp))
}

func (h *HandleUser) forget(c *gin.Context) {
	result := common.Result{}
	//绑定参数
	var req modelUser.ForgetReceiveStruct
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
	msg := &user.ForgetRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		zap.L().Error("evn_api api user forget copy ForgetRequest error")
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "系统内部错误"))
		return
	}

	rsp, err := rpc.UserServiceClient.Forget(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (h *HandleUser) getSpaceIndividual(c *gin.Context) {
	result := common.Result{}
	//绑定参数
	var req modelUser.GetSpaceIndividualReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &user.SpaceIndividualRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		zap.L().Error("evn_api api user getSpaceIndividual copy SpaceIndividualRequest error")
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "系统内部错误"))
		return
	}
	uid := c.GetUint("uid")
	msg.Uid = uint32(uid)
	rsp, err := rpc.UserServiceClient.GetSpaceIndividual(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 SpaceIndividual  实例
	var spaceIndividualJson modelUser.GetSpaceIndividualResponseStruct
	// 将 JSON 字符串解码到 rotographJson  videoListJson实例
	err = json.Unmarshal([]byte(rsp.Data), &spaceIndividualJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(spaceIndividualJson))
}

func (h *HandleUser) getReleaseInformation(c *gin.Context) {
	result := common.Result{}
	//绑定参数
	var req modelUser.GetReleaseInformationReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &user.ReleaseInformationRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		zap.L().Error("evn_api api user getReleaseInformation copy ReleaseInformationResponse error")
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "系统内部错误"))
		return
	}
	rsp, err := rpc.UserServiceClient.GetReleaseInformation(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 videoListJson  实例
	var rspJson modelUser.GetReleaseInformationResponseStruct
	// 将 JSON 字符串解码到rspJson实例
	err = json.Unmarshal([]byte(rsp.Data), &rspJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.Success(rspJson))
}

func (h *HandleUser) getAttentionList(c *gin.Context) {
	result := common.Result{}
	//绑定参数
	var req modelUser.GetAttentionListReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &user.AttentionListRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		zap.L().Error("evn_api api user getAttentionList copy AttentionListRequest error")
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "系统内部错误"))
		return
	}
	rsp, err := rpc.UserServiceClient.GetAttentionList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 attentionListJson  实例
	var attentionListJson modelUser.GetAttentionListInfoList
	// 将 JSON 字符串解码到 attentionListJson实例
	err = json.Unmarshal([]byte(rsp.Data), &attentionListJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(attentionListJson))
}

func (h *HandleUser) getVermicelliList(c *gin.Context) {
	result := common.Result{}
	//绑定参数
	var req modelUser.GetVermicelliListReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &user.VermicelliListRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		zap.L().Error("evn_api api user getVermicelliList copy VermicelliListRequest error")
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "系统内部错误"))
		return
	}
	rsp, err := rpc.UserServiceClient.GetVermicelliList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 vermicelliListJson  实例
	var vermicelliListJson modelUser.GetVermicelliListInfoList
	// 将 JSON 字符串解码到 vermicelliListJson实例
	err = json.Unmarshal([]byte(rsp.Data), &vermicelliListJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(vermicelliListJson))
}
