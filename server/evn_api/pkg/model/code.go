package model

import "dragonsss.cn/evn_common/errs"

// 自定义错误码
var (
	NoLegalMobile = errs.NewError(2001, "手机号不合法")
)
