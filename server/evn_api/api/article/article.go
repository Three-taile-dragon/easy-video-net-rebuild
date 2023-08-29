package article

import (
	"context"
	"dragonsss.cn/evn_api/api/article/rpc"
	articleModel "dragonsss.cn/evn_api/pkg/model/article"
	common "dragonsss.cn/evn_common"
	"dragonsss.cn/evn_common/errs"
	"dragonsss.cn/evn_grpc/article"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"time"
)

type HandleArticle struct {
}

func New() *HandleArticle {
	return &HandleArticle{}
}

// getArticleContributionList 首页查询专栏
func (a *HandleArticle) getArticleContributionList(ctx *gin.Context) {
	result := common.Result{}
	var req articleModel.GetArticleContributionListReceiveStruct
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	c, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()

	msg := &article.CommonPageInfo{
		Page: int32(req.PageInfo.Page),
		Size: int32(req.PageInfo.Size),
	}
	rsp, err := rpc.ArticleServiceClient.GetArticleContributionList(c, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 rspJson  实例
	var rspJson articleModel.GetArticleContributionListByUserResponseList
	//如果没有返回数据
	if rsp.Data == "" {
		ctx.JSON(http.StatusOK, result.Success(&articleModel.GetArticleContributionListByUserResponseList{}))
		return
	}
	// 将 JSON 字符串解码到 GetArticleContributionListByUserResponseList 实例
	err = json.Unmarshal([]byte(rsp.Data), &rspJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	ctx.JSON(http.StatusOK, result.Success(rspJson))
}

// getArticleContributionListByUser 查询用户发布的专栏
func (a *HandleArticle) getArticleContributionListByUser(ctx *gin.Context) {
	result := common.Result{}
	var req articleModel.GetArticleContributionListByUserReceiveStruct
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	c, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()

	msg := &article.CommonIDRequest{
		ID: uint32(req.UserID),
	}
	rsp, err := rpc.ArticleServiceClient.GetArticleContributionListByUser(c, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 rspJson  实例
	var rspJson articleModel.GetArticleContributionListByUserResponseList
	//如果没有返回数据
	if rsp.Data == "" {
		ctx.JSON(http.StatusOK, result.Success(&articleModel.GetArticleContributionListByUserResponseList{}))
		return
	}
	// 将 JSON 字符串解码到 GetArticleContributionListByUserResponseList 实例
	err = json.Unmarshal([]byte(rsp.Data), &rspJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	ctx.JSON(http.StatusOK, result.Success(rspJson))
}

// getArticleComment 获取文章评论
func (a *HandleArticle) getArticleComment(ctx *gin.Context) {
	result := common.Result{}
	var req articleModel.GetArticleCommentReceiveStruct
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	c, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()

	msg := &article.GetArticleCommentRequest{
		PageInfo: &article.CommonPageInfo{
			Page: int32(req.PageInfo.Page),
			Size: int32(req.PageInfo.Size),
		},
		ArticleID: uint32(req.ArticleID),
	}
	rsp, err := rpc.ArticleServiceClient.GetArticleComment(c, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 rspJson  实例
	var rspJson articleModel.GetArticleContributionCommentsResponseStruct
	//如果没有返回数据
	if rsp.Data == "" {
		ctx.JSON(http.StatusOK, result.Success(&articleModel.GetArticleContributionCommentsResponseStruct{}))
		return
	}
	// 将 JSON 字符串解码到 GetArticleContributionCommentsResponseStruct 实例
	err = json.Unmarshal([]byte(rsp.Data), &rspJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	ctx.JSON(http.StatusOK, result.Success(rspJson))
}

// getArticleClassificationList 获取专栏分类
func (a *HandleArticle) getArticleClassificationList(ctx *gin.Context) {
	result := common.Result{}
	//调用grpc
	//对grpc进行两秒超时处理
	c, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()

	rsp, err := rpc.ArticleServiceClient.GetArticleClassificationList(c, &article.CommonZeroRequest{})
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 rspJson  实例
	var rspJson articleModel.ArticleClassificationInfoList
	//如果没有返回数据
	if rsp.Data == "" {
		ctx.JSON(http.StatusOK, result.Success(&articleModel.ArticleClassificationInfoList{}))
		return
	}
	// 将 JSON 字符串解码到 ArticleClassificationInfoList 实例
	err = json.Unmarshal([]byte(rsp.Data), &rspJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	ctx.JSON(http.StatusOK, result.Success(rspJson))
}

// getArticleTotalInfo 获取文章相关总和信息
func (a *HandleArticle) getArticleTotalInfo(ctx *gin.Context) {
	result := common.Result{}
	//调用grpc
	//对grpc进行两秒超时处理
	c, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()

	rsp, err := rpc.ArticleServiceClient.GetArticleTotalInfo(c, &article.CommonZeroRequest{})
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 rspJson  实例
	var rspJson articleModel.GetArticleTotalInfoResponseStruct
	//如果没有返回数据
	if rsp.Data == "" {
		ctx.JSON(http.StatusOK, result.Success(&articleModel.GetArticleTotalInfoResponseStruct{}))
		return
	}
	// 将 JSON 字符串解码到 GetArticleTotalInfoResponseStruct 实例
	err = json.Unmarshal([]byte(rsp.Data), &rspJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	ctx.JSON(http.StatusOK, result.Success(rspJson))
}

// getArticleContributionByID 查询专栏信息根据ID
func (a *HandleArticle) getArticleContributionByID(ctx *gin.Context) {
	result := common.Result{}
	var req articleModel.GetArticleContributionByIDReceiveStruct
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	c, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	uid := ctx.GetInt64("currentUid")
	msg := &article.CommonIDAndUIDRequest{
		ID:  uint32(req.ArticleID),
		UID: uint32(uid),
	}
	rsp, err := rpc.ArticleServiceClient.GetArticleContributionByID(c, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 rspJson  实例
	var rspJson articleModel.GetArticleContributionByIDResponseStruct
	//如果没有返回数据
	if rsp.Data == "" {
		ctx.JSON(http.StatusOK, result.Success(&articleModel.GetArticleContributionByIDResponseStruct{}))
		return
	}
	// 将 JSON 字符串解码到 GetArticleContributionByIDResponseStruct 实例
	err = json.Unmarshal([]byte(rsp.Data), &rspJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	ctx.JSON(http.StatusOK, result.Success(rspJson))
}

// createArticleContribution 发布专栏
func (a *HandleArticle) createArticleContribution(ctx *gin.Context) {
	result := common.Result{}
	var req articleModel.CreateArticleContributionReceiveStruct
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	c, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	uid := ctx.GetInt64("currentUid")
	msg := &article.CreateArticleContributionRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	msg.Uid = uint32(uid)
	rsp, err := rpc.ArticleServiceClient.CreateArticleContribution(c, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	ctx.JSON(http.StatusOK, result.Success(rsp.Data))
}

// updateArticleContribution 更新专栏
func (a *HandleArticle) updateArticleContribution(ctx *gin.Context) {
	result := common.Result{}
	var req articleModel.UpdateArticleContributionReceiveStruct
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	c, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	uid := ctx.GetInt64("currentUid")
	msg := &article.UpdateArticleContributionRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	msg.Uid = uint32(uid)
	rsp, err := rpc.ArticleServiceClient.UpdateArticleContribution(c, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	ctx.JSON(http.StatusOK, result.Success(rsp.Data))
}

// deleteArticleByID 删除专栏
func (a *HandleArticle) deleteArticleByID(ctx *gin.Context) {
	result := common.Result{}
	var req articleModel.DeleteArticleByIDReceiveStruct
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	c, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	uid := ctx.GetInt64("currentUid")
	msg := &article.CommonIDAndUIDRequest{
		ID:  uint32(req.ID),
		UID: uint32(uid),
	}
	rsp, err := rpc.ArticleServiceClient.DeleteArticleByID(c, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	ctx.JSON(http.StatusOK, result.Success(rsp.Data))
}

// articlePostComment 发布评论
func (a *HandleArticle) articlePostComment(ctx *gin.Context) {
	result := common.Result{}
	var req articleModel.ArticlesPostCommentReceiveStruct
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	c, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	uid := ctx.GetInt64("currentUid")
	msg := &article.ArticlePostCommentRequest{
		ArticleID: uint32(req.ArticleID),
		Content:   req.Content,
		ContentID: uint32(req.ContentID),
		Uid:       uint32(uid),
	}
	rsp, err := rpc.ArticleServiceClient.ArticlePostComment(c, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	ctx.JSON(http.StatusOK, result.Success(rsp.Data))
}

// getArticleManagementList 创作中心获取专栏稿件列表
func (a *HandleArticle) getArticleManagementList(ctx *gin.Context) {
	result := common.Result{}
	var req articleModel.GetArticleManagementListReceiveStruct
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//调用grpc
	//对grpc进行两秒超时处理
	c, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	uid := ctx.GetInt64("currentUid")
	msg := &article.GetArticleManagementListRequest{
		PageInfo: &article.CommonPageInfo{
			Page: int32(req.PageInfo.Page),
			Size: int32(req.PageInfo.Size),
		},
		Uid: uint32(uid),
	}
	rsp, err := rpc.ArticleServiceClient.GetArticleManagementList(c, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 rspJson  实例
	var rspJson articleModel.GetArticleManagementListResponseStruct
	//如果没有返回数据
	if rsp.Data == "" {
		ctx.JSON(http.StatusOK, result.Success(&articleModel.GetArticleManagementListResponseStruct{}))
		return
	}
	// 将 JSON 字符串解码到 GetArticleManagementListResponseStruct 实例
	err = json.Unmarshal([]byte(rsp.Data), &rspJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	ctx.JSON(http.StatusOK, result.Success(rspJson))
}
