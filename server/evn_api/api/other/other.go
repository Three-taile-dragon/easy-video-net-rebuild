package other

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

type HandleProject struct {
}

func New() *HandleProject {
	return &HandleProject{}
}
func (p *HandleProject) getHomeInfo(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var req struct {
		PageInfo struct {
			Page int64 `json:"page"`
			Size int64 `json:"size"`
		} `json:"page_info"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	msg := &other.HomeRequest{Page: req.PageInfo.Page, Size: req.PageInfo.Size}
	//msg := &other.HomeRequest{Page: 1, Size: 15}
	rsp, err := rpc.OtherServiceClient.GetHomeInfo(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	// 创建新的 rotographJson videoListJson 实例
	var rotographJson other2.RotographInfoList
	var videoListJson other2.VideoInfoList
	// 将 JSON 字符串解码到 rotographJson  videoListJson实例
	err = json.Unmarshal([]byte(rsp.Rotograph), &rotographJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	err = json.Unmarshal([]byte(rsp.VideoList), &videoListJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	c.JSON(http.StatusOK, result.Success(gin.H{
		"rotograph": rotographJson,
		"videoList": videoListJson,
	}))
}
