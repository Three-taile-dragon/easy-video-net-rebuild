package user

import (
	common "dragonsss.cn/evn_common"
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
