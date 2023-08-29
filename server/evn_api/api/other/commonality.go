package other

import (
	"context"
	"dragonsss.cn/evn_api/api/other/rpc"
	other2 "dragonsss.cn/evn_api/pkg/model/other"
	"dragonsss.cn/evn_api/pkg/service"
	common "dragonsss.cn/evn_common"
	"dragonsss.cn/evn_common/errs"
	"dragonsss.cn/evn_common/model"
	"dragonsss.cn/evn_grpc/other"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HandleCommonality struct {
}

func NewCommonality() *HandleCommonality {
	return &HandleCommonality{}
}

func (c HandleCommonality) upload(ctx *gin.Context) {
	result := common.Result{}
	file, _ := ctx.FormFile("file")
	results, err := service.Upload(file, ctx)
	if err != nil {
		code, msg := errs.ParseGrpcError(model.SystemError)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(results))
}

func (c HandleCommonality) uploadSlice(ctx *gin.Context) {
	result := common.Result{}
	file, _ := ctx.FormFile("file")
	results, err := service.UploadSlice(file, ctx)
	if err != nil {
		code, msg := errs.ParseGrpcError(model.SystemError)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(results))
}

func (c HandleCommonality) uploadCheck(ctx *gin.Context) {
	result := common.Result{}
	var req other2.UploadCheckStruct
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	results, err := service.UploadCheck(&req)
	ctx.JSON(http.StatusOK, result.Success(results))
}

func (c HandleCommonality) uploadMerge(ctx *gin.Context) {
	result := common.Result{}
	var req other2.UploadMergeStruct
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	results, err := service.UploadMerge(&req)
	if err != nil {
		code, msg := errs.ParseGrpcError(model.SystemError)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
	}
	ctx.JSON(http.StatusOK, result.Success(results))
}

func (c HandleCommonality) uploadingMethod(ctx *gin.Context) {
	result := common.Result{}
	var req other2.UploadingMethodStruct
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	context1, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &other.UploadingMethodRequest{
		Method: req.Method,
	}
	rsp, err := rpc.OtherServiceClient.UploadingMethod(context1, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	//4.返回结果
	ctx.JSON(http.StatusOK, result.Success(gin.H{
		"type": rsp.Tp,
	}))
}

func (c HandleCommonality) uploadingDir(ctx *gin.Context) {
	result := common.Result{}
	var req other2.UploadingDirStruct
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	context1, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &other.UploadingDirRequest{
		Interface: req.Interface,
	}
	rsp, err := rpc.OtherServiceClient.UploadingDir(context1, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	//4.返回结果
	ctx.JSON(http.StatusOK, result.Success(gin.H{
		"path":    rsp.Path,
		"quality": rsp.Quality,
	}))
}

func (c HandleCommonality) getFullPathOfImage(ctx *gin.Context) {
	result := common.Result{}
	var req other2.GetFullPathOfImageMethodStruct
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	context1, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &other.GetFullPathOfImageRequest{
		Path: req.Path,
		Type: req.Type,
	}
	rsp, err := rpc.OtherServiceClient.GetFullPathOfImage(context1, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	//4.返回结果
	ctx.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (c HandleCommonality) uploadOss(ctx *gin.Context) {
	result := common.Result{}
	file, _ := ctx.FormFile("file")
	interface1 := ctx.GetString("interface")
	results, err := service.UploadOss(file, interface1, ctx)
	if err != nil {
		code, msg := errs.ParseGrpcError(model.SystemError)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(results))
}

func (c HandleCommonality) search(ctx *gin.Context) {
	result := common.Result{}
	var req other2.SearchStruct
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	uid := ctx.GetInt64("currentUid")
	//调用grpc
	//对grpc进行两秒超时处理
	context1, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &other.SearchRequest{
		Page:    int32(req.PageInfo.Page),
		Size:    int32(req.PageInfo.Size),
		Keyword: req.PageInfo.Keyword,
		Type:    req.Type,
		Uid:     uint32(uid),
	}
	rsp, err := rpc.OtherServiceClient.Search(context1, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	if req.Type == "video" {
		// 创建新的 videoList  实例
		var videoList other2.VideoInfoList
		//如果没有返回数据
		if rsp.Data == "" {
			ctx.JSON(http.StatusOK, result.Success(&other2.VideoInfoList{}))
			return
		}
		// 将 JSON 字符串解码到 VideoInfoList 实例
		err = json.Unmarshal([]byte(rsp.Data), &videoList)
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			ctx.JSON(http.StatusOK, result.Fail(code, msg))
			return
		}

		//4.返回结果
		ctx.JSON(http.StatusOK, result.Success(videoList))
		return
	} else if req.Type == "user" {
		// 创建新的 userList  实例
		var userList other2.UserInfoList
		//如果没有返回数据
		if rsp.Data == "" {
			ctx.JSON(http.StatusOK, result.Success(&other2.UserInfoList{}))
			return
		}
		// 将 JSON 字符串解码到 UserInfoList 实例
		err = json.Unmarshal([]byte(rsp.Data), &userList)
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			ctx.JSON(http.StatusOK, result.Fail(code, msg))
			return
		}

		//4.返回结果
		ctx.JSON(http.StatusOK, result.Success(userList))
		return
	}
	//请求不合法
	code, msg1 := errs.ParseGrpcError(model.RequestError)
	ctx.JSON(http.StatusOK, result.Fail(code, msg1))
	return
}
