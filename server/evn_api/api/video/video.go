package project

import (
	"context"
	"dragonsss.cn/evn_api/api/video/rpc"
	"dragonsss.cn/evn_api/pkg/model/video"
	common "dragonsss.cn/evn_common"
	"dragonsss.cn/evn_common/errs"
	video2 "dragonsss.cn/evn_grpc/video"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type HandleVideo struct {
}

func New() *HandleVideo {
	return &HandleVideo{}
}

func (v HandleVideo) getVideoBarrage(c *gin.Context) {
	result := common.Result{}
	var req video.GetVideoBarrageReceiveStruct
	req.ID = c.Query("id")
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()

	videoID, _ := strconv.ParseUint(req.ID, 0, 19)
	msg := &video2.CommonIDRequest{
		ID: uint32(videoID),
	}
	rsp, err := rpc.VideoServiceClient.GetVideoBarrage(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 videoBarrageJson  实例
	var videoBarrageJson video.BarragesJson
	//如果没有返回数据
	if rsp.Data == "" {
		c.JSON(http.StatusOK, result.BarrageSuccess(c, &video.BarragesJson{}))
		return
	}
	// 将 JSON 字符串解码到 BarragesJson 实例
	err = json.Unmarshal([]byte(rsp.Data), &videoBarrageJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.BarrageSuccess(c, videoBarrageJson))
}

func (v HandleVideo) getVideoBarrageList(c *gin.Context) {
	result := common.Result{}
	var req video.GetVideoBarrageListReceiveStruct
	req.ID = c.Query("id")
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()

	videoID, _ := strconv.ParseUint(req.ID, 0, 19)
	msg := &video2.CommonIDRequest{
		ID: uint32(videoID),
	}
	rsp, err := rpc.VideoServiceClient.GetVideoBarrageList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 videoBarrageInfoJson  实例
	var videoBarrageInfoJson video.BarrageInfoList
	//如果没有返回数据
	if rsp.Data == "" {
		c.JSON(http.StatusOK, result.BarrageSuccess(c, &video.BarrageInfoList{}))
		return
	}
	// 将 JSON 字符串解码到 BarrageInfoList 实例
	err = json.Unmarshal([]byte(rsp.Data), &videoBarrageInfoJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.BarrageSuccess(c, videoBarrageInfoJson))
}

func (v HandleVideo) getVideoComment(c *gin.Context) {
	result := common.Result{}
	var req video.GetVideoCommentReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()

	msg := &video2.GetVideoCommentRequest{
		PageInfo: &video2.CommonPageInfo{
			Page: int32(req.PageInfo.Page),
			Size: int32(req.PageInfo.Size),
		},
		VideoID: uint32(req.VideoID),
	}
	rsp, err := rpc.VideoServiceClient.GetVideoComment(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 videoCommentsJson  实例
	var videoCommentsJson video.GetArticleContributionCommentsResponseStruct
	//如果没有返回数据
	if rsp.Data == "" {
		c.JSON(http.StatusOK, result.BarrageSuccess(c, &video.GetArticleContributionCommentsResponseStruct{}))
		return
	}
	// 将 JSON 字符串解码到 GetArticleContributionCommentsResponseStruct 实例
	err = json.Unmarshal([]byte(rsp.Data), &videoCommentsJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.BarrageSuccess(c, videoCommentsJson))
}

func (v HandleVideo) getVideoContributionByID(c *gin.Context) {
	result := common.Result{}
	var req video.GetVideoContributionByIDReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	uid := c.GetInt64("currentUid")
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()

	msg := &video2.GetVideoContributionByIDRequest{
		VideoID: uint32(req.VideoID),
		Uid:     uint32(uid),
	}
	rsp, err := rpc.VideoServiceClient.GetVideoContributionByID(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 rspJson  实例
	var rspJson video.Response
	//如果没有返回数据
	if rsp.Data == "" {
		c.JSON(http.StatusOK, result.BarrageSuccess(c, &video.Response{}))
		return
	}
	// 将 JSON 字符串解码到 Response 实例
	err = json.Unmarshal([]byte(rsp.Data), &rspJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.BarrageSuccess(c, rspJson))
}
