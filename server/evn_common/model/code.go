package model

import (
	"dragonsss.cn/evn_common/errs"
)

// 定义状态码
// 错误码可以根据模块进行划分 这里 1010 代表 user模块
var (
	SystemError         = errs.NewError(00000, "系统内部错误")
	RedisError          = errs.NewError(99900, "redis错误")
	DBError             = errs.NewError(99901, "数据库错误")
	CopyError           = errs.NewError(99902, "结构体复制错误")
	NoLegalEmail        = errs.NewError(10102001, "邮箱不合法") //不合法的手机号
	CaptchaError        = errs.NewError(10102002, "验证码错误")
	CaptchaNoExist      = errs.NewError(10102003, "验证码不存在或已过期")
	EmailExist          = errs.NewError(10102004, "邮箱已存在")
	AccountExist        = errs.NewError(10102005, "账号已存在")
	NameExist           = errs.NewError(10102006, "用户名已存在")
	AccountNoExist      = errs.NewError(10102008, "账号不存在")
	AccountAndPwdError  = errs.NewError(10102009, "账号密码不正确")
	OrganizationNoExist = errs.NewError(10102010, "组织不存在")
	JsonError           = errs.NewError(10102011, "json转换错误")
	NoLoginError        = errs.NewError(101020012, "用户未登陆")
)
