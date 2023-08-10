package project

import (
	"context"
	"dragonsss.cn/evn_api/api/other/rpc"
	other2 "dragonsss.cn/evn_api/pkg/model/other"
	common "dragonsss.cn/evn_common"
	"dragonsss.cn/evn_common/errs"
	"dragonsss.cn/evn_grpc/other"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HandleLive struct {
}

func NewLive() *HandleLive {
	return &HandleLive{}
}

// GetLiveRoom 获取直播房间
func (l HandleLive) getLiveRoom(c *gin.Context) {
	result := common.Result{}
	uid := c.GetInt64("currentUid")
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &other.CommonIDRequest{
		ID: uint32(uid),
	}
	rsp, err := rpc.OtherServiceClient.GetLiveRoom(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(gin.H{
		"address": rsp.Address,
		"key":     rsp.Key,
	}))
}

// GetLiveRoomInfo 获取直播房间信息
func (l HandleLive) getLiveRoomInfo(c *gin.Context) {
	result := common.Result{}
	var req other2.GetLiveRoomInfoReceiveStruct
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
	msg := &other.CommonIDAndUIDRequest{
		ID:  uint32(req.RoomID),
		UID: uint32(uid),
	}
	rsp, err := rpc.OtherServiceClient.GetLiveRoomInfo(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 liveRoomInfoJson  实例
	var liveRoomInfoJson other2.GetLiveRoomInfoResponseStruct
	// 将 JSON 字符串解码到 GetLiveRoomInfoResponseStruct实例
	err = json.Unmarshal([]byte(rsp.Data), &liveRoomInfoJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(liveRoomInfoJson))
}

// GetBeLiveList 获取正在直播的用户
func (l HandleLive) getBeLiveList(c *gin.Context) {
	result := common.Result{}
	uid := c.GetInt64("currentUid")
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &other.CommonIDRequest{
		ID: uint32(uid),
	}
	rsp, err := rpc.OtherServiceClient.GetBeLiveList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 beLiveInfoJson  实例
	var beLiveInfoJson other2.BeLiveInfoList
	// 将 JSON 字符串解码到 BeLiveInfoList实例
	err = json.Unmarshal([]byte(rsp.Data), &beLiveInfoJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(beLiveInfoJson))
}
