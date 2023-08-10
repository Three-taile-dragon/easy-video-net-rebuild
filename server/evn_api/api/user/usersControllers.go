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
	uid := c.GetInt64("currentUid")
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
	uid := c.GetInt64("currentUid")
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
	uid := c.GetInt64("currentUid")
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

func (u *HandleUserControllers) updateAvatar(c *gin.Context) {
	result := &common.Result{}
	//获取传入的邮箱
	//绑定参数
	var req modelUser.UpdateAvatarStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	uid := c.GetInt64("currentUid")
	msg := &user.UpdateAvatarRequest{
		ImgUrl: req.ImgUrl,
		TP:     req.Tp,
		ID:     uid,
	}
	//通过grpc调用 验证码生成函数
	rsp, err := rpc.UserServiceClient.UpdateAvatar(ctx, msg)
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (u *HandleUserControllers) getLiveData(c *gin.Context) {
	result := common.Result{}
	uid := c.GetInt64("currentUid")
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &user.CommonIDRequest{
		ID: uint32(uid),
	}
	rsp, err := rpc.UserServiceClient.GetLiveData(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.Success(gin.H{
		"title": rsp.Title,
		"img":   rsp.Img,
	}))
}

func (u *HandleUserControllers) saveLiveData(c *gin.Context) {
	result := &common.Result{}
	//获取传入的邮箱
	//绑定参数
	var req modelUser.SaveLiveDataReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	uid := c.GetInt64("currentUid")
	msg := &user.SaveLiveDataRequest{
		Img:   req.ImgUrl,
		TP:    req.Tp,
		Title: req.Title,
		ID:    uid,
	}
	//通过grpc调用 验证码生成函数
	rsp, err := rpc.UserServiceClient.SaveLiveData(ctx, msg)
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (u *HandleUserControllers) sendEmailVerificationCodeByChangePassword(c *gin.Context) {
	result := &common.Result{}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	uid := c.GetInt64("currentUid")
	msg := &user.CommonIDRequest{
		ID: uint32(uid),
	}
	//通过grpc调用 验证码生成函数
	rsp, err := rpc.UserServiceClient.SendEmailVerificationCodeByChangePassword(ctx, msg)
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (u *HandleUserControllers) changePassword(c *gin.Context) {
	result := &common.Result{}
	//获取传入的邮箱
	//绑定参数
	var req modelUser.ChangePasswordReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	uid := c.GetInt64("currentUid")
	msg := &user.ChangePasswordRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		zap.L().Error("evn_api user userControllers changePassword Copy err", zap.Error(err))
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "系统错误"))
		return
	}
	msg.ID = uid
	//通过grpc调用
	rsp, err := rpc.UserServiceClient.ChangePassword(ctx, msg)
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (u *HandleUserControllers) attention(c *gin.Context) {
	result := &common.Result{}
	var req modelUser.AttentionReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	curUid := c.GetInt64("currentUid")
	msg := &user.CommonIDAndUIDRequest{
		ID:  uint32(req.Uid), //关注的人的ID
		UID: uint32(curUid),  //用户ID
	}
	//通过grpc调用 验证码生成函数
	rsp, err := rpc.UserServiceClient.Attention(ctx, msg)
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (u *HandleUserControllers) createFavorites(c *gin.Context) {
	result := &common.Result{}
	//获取传入的邮箱
	//绑定参数
	var req modelUser.CreateFavoritesReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &user.FavoritesRequest{}
	err = copier.Copy(msg, req)
	if err != nil {
		zap.L().Error("evn_api user userControllers createFavorites Copy err", zap.Error(err))
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "系统错误"))
		return
	}
	msg.Uid = uint32(c.GetInt64("currentUid"))
	//通过grpc调用
	rsp, err := rpc.UserServiceClient.CreateFavorites(ctx, msg)
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (u *HandleUserControllers) getFavoritesList(c *gin.Context) {
	result := &common.Result{}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	uid := c.GetInt64("currentUid")
	msg := &user.CommonIDRequest{
		ID: uint32(uid),
	}
	//通过grpc调用 验证码生成函数
	rsp, err := rpc.UserServiceClient.GetFavoritesList(ctx, msg)
	// 创建新的 GetFavoritesInfoList  实例
	var list modelUser.GetFavoritesInfoList
	// 将 JSON 字符串解码到 list
	err = json.Unmarshal([]byte(rsp.Data), &list)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(list))
}

func (u *HandleUserControllers) deleteFavorites(c *gin.Context) {
	result := &common.Result{}
	var req modelUser.DeleteFavoritesReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	curUid := c.GetInt64("currentUid")
	msg := &user.CommonIDAndUIDRequest{
		ID:  uint32(req.ID), //关注的人的ID
		UID: uint32(curUid), //用户ID
	}
	//通过grpc调用 验证码生成函数
	rsp, err := rpc.UserServiceClient.DeleteFavorites(ctx, msg)
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (u *HandleUserControllers) favoriteVideo(c *gin.Context) {
	result := &common.Result{}
	var req modelUser.FavoriteVideoReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	curUid := c.GetInt64("currentUid")
	msg := &user.FavoriteVideoRequest{
		IDs:      req.IDs,
		Video_ID: uint32(req.VideoID),
		UID:      uint32(curUid),
	}
	//通过grpc调用 验证码生成函数
	rsp, err := rpc.UserServiceClient.FavoriteVideo(ctx, msg)
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (u *HandleUserControllers) getFavoritesListByFavoriteVideo(c *gin.Context) {
	result := &common.Result{}
	var req modelUser.GetFavoritesListByFavoriteVideoReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	curUid := c.GetInt64("currentUid")
	msg := &user.FavoritesListRequest{
		Video_ID: uint32(req.VideoID),
		UID:      uint32(curUid),
	}
	//通过grpc调用 生成函数
	rsp, err := rpc.UserServiceClient.GetFavoritesListByFavoriteVideo(ctx, msg)
	// 创建新的 GetFavoritesListByFavoriteVideoInfoList  实例
	var flist modelUser.GetFavoritesListByFavoriteVideoInfoList
	// 将 JSON 字符串解码到 list
	err = json.Unmarshal([]byte(rsp.Data), &flist)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(flist))
}

func (u *HandleUserControllers) getFavoriteVideoList(c *gin.Context) {
	result := &common.Result{}
	var req modelUser.GetFavoriteVideoListReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &user.FavoriteVideoListRequest{
		Favorite_ID: uint32(req.FavoriteID),
	}
	//通过grpc调用函数
	rsp, err := rpc.UserServiceClient.GetFavoriteVideoList(ctx, msg)
	// 创建新的 GetFavoriteVideoListResponseStruct  实例
	var flrs modelUser.GetFavoriteVideoListResponseStruct
	// 将 JSON 字符串解码到 flrs
	err = json.Unmarshal([]byte(rsp.Data), &flrs)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(flrs))
}

func (u *HandleUserControllers) getRecordList(c *gin.Context) {
	result := &common.Result{}
	var req modelUser.GetRecordListReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	curUid := c.GetInt64("currentUid")
	msg := &user.GetRecordListRequest{
		Page: int64(req.PageInfo.Page),
		Size: int64(req.PageInfo.Size),
		Uid:  uint32(curUid),
	}
	//通过grpc调用 生成函数
	rsp, err := rpc.UserServiceClient.GetRecordList(ctx, msg)
	// 创建新的 GetRecordListItemList  实例
	var rl modelUser.GetRecordListItemList
	// 将 JSON 字符串解码到 list
	err = json.Unmarshal([]byte(rsp.Data), &rl)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(rl))
}

func (u *HandleUserControllers) clearRecord(c *gin.Context) {
	result := common.Result{}
	uid := c.GetInt64("currentUid")
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &user.CommonIDRequest{
		ID: uint32(uid),
	}
	rsp, err := rpc.UserServiceClient.ClearRecord(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (u *HandleUserControllers) deleteRecordByID(c *gin.Context) {
	result := &common.Result{}
	var req modelUser.DeleteRecordByIDReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	curUid := c.GetInt64("currentUid")
	msg := &user.CommonIDAndUIDRequest{
		ID:  uint32(req.ID), //指定历史记录的ID
		UID: uint32(curUid), //用户ID
	}
	//通过grpc调用 验证码生成函数
	rsp, err := rpc.UserServiceClient.DeleteRecordByID(ctx, msg)
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (u *HandleUserControllers) getNoticeList(c *gin.Context) {
	result := &common.Result{}
	var req modelUser.GetNoticeListReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	curUid := c.GetInt64("currentUid")
	msg := &user.GetNoticeListRequest{
		Tp:   req.Type,
		Page: int64(req.PageInfo.Page),
		Size: int64(req.PageInfo.Size),
		Uid:  uint32(curUid),
	}
	//通过grpc调用 验证码生成函数
	rsp, err := rpc.UserServiceClient.GetNoticeList(ctx, msg)
	// 创建新的 GetNoticeListStruct  实例
	var nls modelUser.GetNoticeListStruct
	// 将 JSON 字符串解码到 list
	err = json.Unmarshal([]byte(rsp.Data), &nls)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(nls))
}

func (u *HandleUserControllers) getChatList(c *gin.Context) {
	result := common.Result{}
	uid := c.GetInt64("currentUid")
	//调用grpc
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &user.CommonIDRequest{
		ID: uint32(uid),
	}
	rsp, err := rpc.UserServiceClient.GetChatList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误信息
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}

	// 创建新的 chatListJson  实例
	var chatListJson modelUser.GetChatListResponseStruct
	// 将 JSON 字符串解码到 GetChatListResponseStruct
	err = json.Unmarshal([]byte(rsp.Data), &chatListJson)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(chatListJson))
}

func (u *HandleUserControllers) getChatHistoryMsg(c *gin.Context) {
	result := &common.Result{}
	var req modelUser.GetChatHistoryMsgStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	curUid := c.GetInt64("currentUid")
	msg := &user.GetChatHistoryMsgRequest{
		Tid:      uint32(req.Tid),
		LastTime: req.LastTime.String(),
		Uid:      uint32(curUid),
	}
	//通过grpc调用 生成函数
	rsp, err := rpc.UserServiceClient.GetChatHistoryMsg(ctx, msg)
	// 创建新的 []ChatMessageInfo  实例
	var cmi []modelUser.ChatMessageInfo
	// 将 JSON 字符串解码到 list
	err = json.Unmarshal([]byte(rsp.Data), &cmi)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	//4.返回结果
	c.JSON(http.StatusOK, result.Success(cmi))
}

func (u *HandleUserControllers) personalLetter(c *gin.Context) {
	result := &common.Result{}
	var req modelUser.PersonalLetterReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	curUid := c.GetInt64("currentUid")
	msg := &user.CommonIDAndUIDRequest{
		ID:  uint32(req.ID),
		UID: uint32(curUid), //用户ID
	}
	//通过grpc调用 函数
	rsp, err := rpc.UserServiceClient.PersonalLetter(ctx, msg)
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}

func (u *HandleUserControllers) deleteChatItem(c *gin.Context) {
	result := &common.Result{}
	var req modelUser.DeleteChatItemReceiveStruct
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	//对grpc进行两秒超时处理
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	curUid := c.GetInt64("currentUid")
	msg := &user.CommonIDAndUIDRequest{
		ID:  uint32(req.ID),
		UID: uint32(curUid), //用户ID
	}
	//通过grpc调用 函数
	rsp, err := rpc.UserServiceClient.DeleteChatItem(ctx, msg)
	//结果返回
	if err != nil {
		code, msg := errs.ParseGrpcError(err) //解析grpc错误
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp.Data))
}
