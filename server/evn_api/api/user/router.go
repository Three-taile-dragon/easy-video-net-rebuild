package user

import (
	"dragonsss.cn/evn_api/api/cors"
	"dragonsss.cn/evn_api/api/midd"
	"dragonsss.cn/evn_api/api/user/rpc"
	"dragonsss.cn/evn_api/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

type RouterUser struct {
}

func init() {
	log.Println("init user router")
	zap.L().Info("init user router")
	ru := &RouterUser{}
	router.Register(ru)
}

func (*RouterUser) Router(r *gin.Engine) {
	//初始化grpc客户端连接
	//使得可以调用逻辑函数
	rpc.InitRpcUserClient()
	h := New()
	group := r.Group("/api/user")
	group.Use(cors.Cors())
	group.POST("/getCaptcha", h.getCaptcha)
	group.POST("/login", h.login)
	group.POST("/register", h.register)
	group.POST("/refreshToken", h.refreshToken)
	group.POST("/forget", h.forget)
	group.POST("/space/getSpaceIndividual", h.getSpaceIndividual)
	group.POST("/space/getReleaseInformation", h.getReleaseInformation)
	//必须登入
	spaceRouter := r.Group("/api/user/space").Use(midd.TokenVerify())
	{
		spaceRouter.POST("/getAttentionList", h.getAttentionList)
		spaceRouter.POST("/getVermicelliList", h.getVermicelliList)
	}
	userRouter := r.Group("/api/user").Use(midd.TokenVerify())
	{
		u := NewUserControllers()
		userRouter.POST("/getUserInfo", u.getUserInfo)
		userRouter.POST("/setUserInfo", u.setUserInfo)
		userRouter.POST("/determineNameExists", u.determineNameExists)
		userRouter.POST("/updateAvatar", u.updateAvatar)
		userRouter.POST("/getLiveData", u.getLiveData)
		userRouter.POST("/saveLiveData", u.saveLiveData)
		userRouter.POST("/sendEmailVerificationCodeByChangePassword", u.sendEmailVerificationCodeByChangePassword)
		userRouter.POST("/changePassword", u.changePassword)
		userRouter.POST("/attention", u.attention)
		userRouter.POST("/createFavorites", u.createFavorites)
		userRouter.POST("/getFavoritesList", u.getFavoritesList)
		userRouter.POST("/deleteFavorites", u.deleteFavorites)
		userRouter.POST("/favoriteVideo", u.favoriteVideo)
		userRouter.POST("/getFavoritesListByFavoriteVideo", u.getFavoritesListByFavoriteVideo)
		userRouter.POST("/getFavoriteVideoList", u.getFavoriteVideoList)
		userRouter.POST("/getRecordList", u.getRecordList)
		userRouter.POST("/clearRecord", u.clearRecord)
		userRouter.POST("/deleteRecordByID", u.deleteRecordByID)
		userRouter.POST("/getNoticeList", u.getNoticeList)
		userRouter.POST("/getChatList", u.getChatList)
		userRouter.POST("/getChatHistoryMsg", u.getChatHistoryMsg)
		userRouter.POST("/personalLetter", u.personalLetter)
		userRouter.POST("/deleteChatItem", u.deleteChatItem)
	}
}
