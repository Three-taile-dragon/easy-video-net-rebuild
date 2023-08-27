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
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/copier"
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
	c.JSON(http.StatusOK, result.Success(videoCommentsJson))
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
		c.JSON(http.StatusOK, result.Success(&video.Response{}))
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
	c.JSON(http.StatusOK, result.Success(rspJson))
}

func (v HandleVideo) sendVideoBarrage(c *gin.Context) {
	result := common.Result{}
	var req video.SendVideoBarrageReceiveStruct
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	uid := c.GetInt64("currentUid")
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()

	msg := &video2.SendVideoBarrageRequest{
		Author: req.Author,
		Color:  uint32(req.Color),
		ID:     req.ID,
		Text:   req.Text,
		Time:   float32(req.Time),
		Type:   uint32(req.Type),
		Uid:    uint32(uid),
	}
	rsp, err := rpc.VideoServiceClient.SendVideoBarrage(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 rspJson  实例
	var rspJson video.SendVideoBarrageReceiveStruct
	//如果没有返回数据
	if rsp.Data == "" {
		c.JSON(http.StatusOK, result.BarrageSuccess(c, &video.SendVideoBarrageReceiveStruct{}))
		return
	}
	// 将 JSON 字符串解码到 SendVideoBarrageReceiveStruct 实例
	err = json.Unmarshal([]byte(rsp.Data), &rspJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.BarrageSuccess(c, rspJson))
}

func (v HandleVideo) createVideoContribution(c *gin.Context) {
	result := common.Result{}
	var req video.CreateVideoContributionReceiveStruct
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	uid := c.GetInt64("currentUid")
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()

	msg := &video2.CreateVideoContributionRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	msg.Uid = uint32(uid)
	//msg.Reprinted = *req.Reprinted
	//msg.Media = *req.Media
	rsp, err := rpc.VideoServiceClient.CreateVideoContribution(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (v HandleVideo) updateVideoContribution(c *gin.Context) {
	result := common.Result{}
	var req video.UpdateVideoContributionReceiveStruct
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

	msg := &video2.UpdateVideoContributionRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	msg.Uid = uint32(uid)
	rsp, err := rpc.VideoServiceClient.UpdateVideoContribution(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (v HandleVideo) deleteVideoByID(c *gin.Context) {
	result := common.Result{}
	var req video.DeleteVideoByIDReceiveStruct
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

	msg := &video2.CommonIDAndUIDRequest{
		ID:  uint32(req.ID),
		UID: uint32(uid),
	}
	rsp, err := rpc.VideoServiceClient.DeleteVideoByID(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (v HandleVideo) videoPostComment(c *gin.Context) {
	result := common.Result{}
	var req video.VideosPostCommentReceiveStruct
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

	msg := &video2.VideoPostCommentRequest{
		VideoID:   uint32(req.VideoID),
		Content:   req.Content,
		ContentID: uint32(req.ContentID),
		Uid:       uint32(uid),
	}
	rsp, err := rpc.VideoServiceClient.VideoPostComment(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (v HandleVideo) getVideoManagementList(c *gin.Context) {
	result := common.Result{}
	var req video.GetVideoManagementListReceiveStruct
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

	msg := &video2.GetVideoManagementListRequest{
		PageInfo: &video2.CommonPageInfo{
			Page: int32(req.PageInfo.Page),
			Size: int32(req.PageInfo.Size),
		},
		Uid: uint32(uid),
	}
	rsp, err := rpc.VideoServiceClient.GetVideoManagementList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	// 创建新的 rspJson  实例
	var rspJson video.GetVideoManagementList
	//如果没有返回数据
	if rsp.Data == "" {
		c.JSON(http.StatusOK, result.Success(&video.GetVideoManagementList{}))
		return
	}
	// 将 JSON 字符串解码到 GetVideoManagementList 实例
	err = json.Unmarshal([]byte(rsp.Data), &rspJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.Success(rspJson))
}

func (v HandleVideo) likeVideo(c *gin.Context) {
	result := common.Result{}
	var req video.LikeVideoReceiveStruct
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

	msg := &video2.CommonIDAndUIDRequest{
		ID:  uint32(req.VideoID),
		UID: uint32(uid),
	}
	rsp, err := rpc.VideoServiceClient.LikeVideo(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}
