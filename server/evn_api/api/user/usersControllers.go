package user

import (
	"context"
	"dragonsss.cn/evn_api/api/user/rpc"
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

type HandleUserControllers struct {
}

func NewUserControllers() *HandleUserControllers {
	return &HandleUserControllers{}
}

func (u *HandleUserControllers) getUserInfo(c *gin.Context) {
	result := common.Result{}
	uid := c.GetInt64("uid")
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &user.CommonIDRequest{
		ID: uint32(uid),
	}
	rsp, err := rpc.UserServiceClient.GetUserInfo(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 userSetInfoJson  实例
	var userSetInfoJson modelUser.UserSetInfoResponseStruct
	// 将 JSON 字符串解码到 vermicelliListJson实例
	err = json.Unmarshal([]byte(rsp.Data), &userSetInfoJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(userSetInfoJson))
}

func (u *HandleUserControllers) setUserInfo(c *gin.Context) {
	result := common.Result{}
	//绑定参数
	var req modelUser.SetUserInfoReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	uid := c.GetInt64("uid")
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &user.UserInfoRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		zap.L().Error("evn_api api userControllers setUserInfo copy UserInfoRequest error", zap.Error(err))
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "系统内部错误"))
		return
	}
	msg.ID = uid
	msg.Birth_Date = req.BirthDate
	rsp, err := rpc.UserServiceClient.SetUserInfo(ctx, msg)
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp))
}

func (u *HandleUserControllers) determineNameExists(c *gin.Context) {
	result := &common.Result{}
	//获取传入的邮箱
	//绑定参数
	var req modelUser.DetermineNameExistsStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &user.DetermineNameExistsRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		zap.L().Error("evn_api api userControllers setUserInfo copy UserInfoRequest error", zap.Error(err))
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "系统内部错误"))
		return
	}
	uid := c.GetInt64("uid")
	msg.ID = uid
	//通过grpc调用 验证码生成函数
	rsp, err := rpc.UserServiceClient.DetermineNameExists(ctx, msg)
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}
