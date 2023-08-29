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

type HandleDiscuss struct {
}

func NewDiscuss() *HandleDiscuss {
	return &HandleDiscuss{}
}

func (d HandleDiscuss) getDiscussVideoList(c *gin.Context) {
	result := common.Result{}
	var req other2.GetDiscussVideoListReceiveStruct
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
	msg := &other.CommonDiscussRequest{
		Page: int32(req.PageInfo.Page),
		Size: int32(req.PageInfo.Size),
		Uid:  uint32(uid),
	}
	rsp, err := rpc.OtherServiceClient.GetDiscussVideoList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 discussVideoListJson  实例
	var discussVideoListJson other2.GetDiscussVideoListStruct
	//如果没有返回数据
	if rsp.Data == "" {
		c.JSON(http.StatusOK, result.Success(&other2.GetDiscussVideoListStruct{}))
		return
	}
	// 将 JSON 字符串解码到 GetDiscussVideoListStruct实例
	err = json.Unmarshal([]byte(rsp.Data), &discussVideoListJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(discussVideoListJson))
}

func (d HandleDiscuss) getDiscussArticleList(c *gin.Context) {
	result := common.Result{}
	var req other2.GetDiscussVideoListReceiveStruct
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
	msg := &other.CommonDiscussRequest{
		Page: int32(req.PageInfo.Page),
		Size: int32(req.PageInfo.Size),
		Uid:  uint32(uid),
	}
	rsp, err := rpc.OtherServiceClient.GetDiscussArticleList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 discussArticleListJson  实例
	var discussArticleListJson other2.GetDiscussArticleListStruct
	//如果没有返回数据
	if rsp.Data == "" {
		c.JSON(http.StatusOK, result.Success(&other2.GetDiscussArticleListStruct{}))
		return
	}
	// 将 JSON 字符串解码到 GetDiscussArticleListStruct 实例
	err = json.Unmarshal([]byte(rsp.Data), &discussArticleListJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(discussArticleListJson))
}

func (d HandleDiscuss) getDiscussBarrageList(c *gin.Context) {
	result := common.Result{}
	var req other2.GetDiscussVideoListReceiveStruct
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
	msg := &other.CommonDiscussRequest{
		Page: int32(req.PageInfo.Page),
		Size: int32(req.PageInfo.Size),
		Uid:  uint32(uid),
	}
	rsp, err := rpc.OtherServiceClient.GetDiscussBarrageList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 discussBarrageListJson  实例
	var discussBarrageListJson other2.GetDiscussBarrageListStruct
	//如果没有返回数据
	if rsp.Data == "" {
		c.JSON(http.StatusOK, result.Success(&other2.GetDiscussBarrageListStruct{}))
		return
	}
	// 将 JSON 字符串解码到 GetDiscussBarrageListStruct 实例
	err = json.Unmarshal([]byte(rsp.Data), &discussBarrageListJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(discussBarrageListJson))
}
