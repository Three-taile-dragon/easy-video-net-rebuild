package user

import (
	common "dragonsss.cn/evn_common"
	common2 "dragonsss.cn/evn_common/model/common"
	"errors"
	"time"
)

// RegisterReq user 相关模型
// 注意加上 form 表单
type RegisterReq struct {
	Email     string `json:"email" form:"email"`
	Name      string `json:"name" form:"name"`
	Password  string `json:"password" form:"password"`
	Password2 string `json:"password2" form:"password2"`
	//Mobile    string `json:"mobile" form:"mobile"`
	Captcha string `json:"captcha" form:"captcha"`
}

func (r RegisterReq) VerifyPassword() bool {
	return r.Password == r.Password2
}

// Verify 验证参数
func (r RegisterReq) Verify() error {
	//验证 邮箱 手机号 密码 用户名等等是否合法
	if !common.VerifyEmailFormat(r.Email) {
		return errors.New("邮箱格式不正确")
	}
	//if !common.VerifyMobile(r.Mobile) {
	//	return errors.New("手机号格式不正确")
	//}
	if !r.VerifyPassword() {
		return errors.New("两次密输入不一致")
	}
	return nil
}

type LoginReq struct {
	Username string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRsp struct {
	ID        uint      `json:"id"`
	UserName  string    `json:"username"`
	Photo     string    `json:"photo"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}
type EmailCaptcha struct {
	Email string `json:"email" binding:"required,email"`
}

type TokenList struct {
	AccessToken    string `json:"accessToken"`
	RefreshToken   string `json:"refreshToken"`
	TokenType      string `json:"tokenType"`
	AccessTokenExp int64  `json:"accessTokenExp"`
}

func (l *LoginReq) VerifyAccount() error {
	if l.Username == "" {
		return errors.New("用户名不能为空")
	}
	return nil
}

// ForgetReceiveStruct 用户找回密码
type ForgetReceiveStruct struct {
	Password  string `json:"password" form:"required"`
	Password2 string `json:"password2" form:"password2"`
	Email     string `json:"email" form:"required,email"`
	Captcha   string `json:"captcha" form:"captcha"`
}

func (f ForgetReceiveStruct) VerifyPassword() bool {
	return f.Password == f.Password2
}

// Verify 验证参数
func (f ForgetReceiveStruct) Verify() error {
	//验证 邮箱 手机号 密码 用户名等等是否合法
	if !common.VerifyEmailFormat(f.Email) {
		return errors.New("邮箱格式不正确")
	}
	//if !common.VerifyMobile(r.Mobile) {
	//	return errors.New("手机号格式不正确")
	//}
	if !f.VerifyPassword() {
		return errors.New("两次密码输入不一致")
	}
	return nil
}

type GetSpaceIndividualReceiveStruct struct {
	ID uint `json:"id" binding:"required"`
}

type GetSpaceIndividualResponseStruct struct {
	ID            uint   `json:"id"`
	UserName      string `json:"username"`
	Photo         string `json:"photo"`
	Signature     string `json:"signature"`
	IsAttention   bool   `json:"is_attention"`
	AttentionNum  *int64 `json:"attention_num"`
	VermicelliNum *int64 `json:"vermicelli_num"`
}

type GetReleaseInformationReceiveStruct struct {
	ID uint `json:"id" binding:"required"`
}

type ReleaseInformationVideoInfo struct {
	ID            uint      `json:"id"`
	Uid           uint      `json:"uid" `
	Title         string    `json:"title" `
	Video         string    `json:"video"`
	Cover         string    `json:"cover" `
	VideoDuration int64     `json:"video_duration"`
	Label         []string  `json:"label"`
	Introduce     string    `json:"introduce"`
	Heat          int       `json:"heat"`
	BarrageNumber int       `json:"barrageNumber"`
	Username      string    `json:"username"`
	CreatedAt     time.Time `json:"created_at"`
}

type ReleaseInformationVideoInfoList []ReleaseInformationVideoInfo

type ReleaseInformationArticleInfo struct {
	Id             uint      `json:"id"`
	Uid            uint      `json:"uid" `
	Title          string    `json:"title" `
	Cover          string    `json:"cover" `
	Label          []string  `json:"label" `
	Content        string    `json:"content"`
	IsComments     bool      `json:"is_comments"`
	Heat           int       `json:"heat"`
	LikesNumber    int       `json:"likes_number"`
	CommentsNumber int       `json:"comments_number"`
	Classification string    `json:"classification"`
	CreatedAt      time.Time `json:"created_at"`
}

type ReleaseInformationArticleInfoList []ReleaseInformationArticleInfo

type GetReleaseInformationResponseStruct struct {
	VideoList   ReleaseInformationVideoInfoList   `json:"videoList"`
	ArticleList ReleaseInformationArticleInfoList `json:"articleList"`
}

type GetAttentionListReceiveStruct struct {
	ID uint `json:"id" binding:"required"`
}

type GetAttentionListInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Signature   string `json:"signature"`
	Photo       string `json:"photo"`
	IsAttention bool   `json:"is_attention"`
}

type GetAttentionListInfoList []GetAttentionListInfo

type GetVermicelliListReceiveStruct struct {
	ID uint `json:"id" binding:"required"`
}

type GetVermicelliListInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Signature   string `json:"signature"`
	Photo       string `json:"photo"`
	IsAttention bool   `json:"is_attention"`
}

type GetVermicelliListInfoList []GetVermicelliListInfo

type UserSetInfoResponseStruct struct {
	ID        uint      `json:"id"`
	UserName  string    `json:"username"`
	Gender    int8      `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
	IsVisible bool      `json:"is_visible"`
	Signature string    `json:"signature"`
}

// SetUserInfoReceiveStruct 设置用户信息
type SetUserInfoReceiveStruct struct {
	Username  string `json:"username" binding:"required"`
	Gender    *int   `json:"gender" binding:"required"`
	BirthDate string `json:"birth_Date" binding:"required"`
	IsVisible *bool  `json:"is_Visible" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type DetermineNameExistsStruct struct {
	Username string `json:"username" binding:"required"`
}

// UpdateAvatarStruct 更新头像
type UpdateAvatarStruct struct {
	ImgUrl string `json:"imgUrl" binding:"required"`
	Tp     string `json:"type" binding:"required"`
}

// SaveLiveDataReceiveStruct 设置直播信息
type SaveLiveDataReceiveStruct struct {
	Tp     string `json:"type" binding:"required"`
	ImgUrl string `json:"imgUrl" binding:"required"`
	Title  string `json:"title" binding:"required"`
}

// ChangePasswordReceiveStruct 修改密码
type ChangePasswordReceiveStruct struct {
	Captcha         string `json:"captcha" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type AttentionReceiveStruct struct {
	Uid uint `json:"uid"  binding:"required" binding:"required"`
}

type CreateFavoritesReceiveStruct struct {
	ID      uint   `json:"id"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	Cover   string `json:"cover"`
	Tp      string `json:"type"`
}

type GetFavoritesInfo struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Cover    string `json:"cover"`
	Tp       string `json:"type"`
	Src      string `json:"src"`
	Max      int    `json:"max"`
	UsesInfo struct {
		Username string `json:"username"`
	} `json:"userInfo"`
}

type GetFavoritesInfoList []GetFavoritesInfo

type DeleteFavoritesReceiveStruct struct {
	ID uint `json:"id" binding:"required"`
}

type FavoriteVideoReceiveStruct struct {
	IDs     []uint32 `json:"ids" binding:"required"`
	VideoID uint     `json:"video_id" binding:"required"`
}

type GetFavoritesListByFavoriteVideoReceiveStruct struct {
	VideoID uint `json:"video_id" binding:"required"`
}

type GetFavoritesListByFavoriteVideoInfo struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Cover    string `json:"cover"`
	Tp       string `json:"type"`
	Src      string `json:"src"`
	Max      int    `json:"max"`
	Selected bool   `json:"selected"`
	Present  int    `json:"present"`
	UsesInfo struct {
		Username string `json:"username"`
	} `json:"userInfo"`
}

type GetFavoritesListByFavoriteVideoInfoList []GetFavoritesListByFavoriteVideoInfo

type GetFavoriteVideoListReceiveStruct struct {
	FavoriteID uint `json:"favorite_id" binding:"required"`
}

type GetFavoriteVideoListItem struct {
	ID            uint      `json:"id"`
	Uid           uint      `json:"uid"`
	Title         string    `json:"title"`
	Video         string    `json:"video"`
	Cover         string    `json:"cover"`
	VideoDuration int64     `json:"video_duration"`
	CreatedAt     time.Time `json:"created_at"`
}

type GetFavoriteVideoList []GetFavoriteVideoListItem

type GetFavoriteVideoListResponseStruct struct {
	VideoList GetFavoriteVideoList `json:"videoList"`
}

type GetRecordListReceiveStruct struct {
	PageInfo common2.PageInfo `json:"page_info" binding:"required"`
}

type GetRecordListItem struct {
	ID        uint      `json:"id"`
	ToID      uint      `json:"to_id"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Username  string    `json:"username"`
	Photo     string    `json:"photo"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetRecordListItemList []GetRecordListItem

type DeleteRecordByIDReceiveStruct struct {
	ID uint `json:"id"`
}

type GetNoticeListReceiveStruct struct {
	Type     string           `json:"type"`
	PageInfo common2.PageInfo `json:"page_info" binding:"required"`
}

type GetNoticeListItem struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Type      string    `json:"type"`
	ToID      uint      `json:"to_id"`
	Photo     string    `json:"photo"`
	Comment   string    `json:"comment"`
	Cover     string    `json:"cover"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type GetNoticeListStruct []GetNoticeListItem

type ChatMessageInfo struct {
	ID        uint      `json:"id"`
	Uid       uint      `json:"uid"`
	Username  string    `json:"username"`
	Photo     string    `json:"photo"`
	Tid       uint      `json:"tid"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type GetChatListItem struct {
	ToID            uint              `json:"to_id"`
	Username        string            `json:"username"`
	Photo           string            `json:"photo"`
	Unread          int               `json:"unread" gorm:"unread"`
	LastMessage     string            `json:"last_message"`
	LastMessagePage int               `json:"last_message_page"`
	MessageList     []ChatMessageInfo `json:"message_list"`
	LastAt          time.Time         `json:"last_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type GetChatListResponseStruct []GetChatListItem

type GetChatHistoryMsgStruct struct {
	Tid      uint      `json:"tid"`
	LastTime time.Time `json:"last_time"`
}

type PersonalLetterReceiveStruct struct {
	ID uint `json:"id" binding:"required"`
}

type DeleteChatItemReceiveStruct struct {
	ID uint `json:"id" binding:"required"`
}
