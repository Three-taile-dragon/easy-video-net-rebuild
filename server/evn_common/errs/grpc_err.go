package errs

import (
	common "dragonsss.cn/evn_common"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GrpcError 处理grpc错误信息
func GrpcError(err *BError) error {
	return status.Error(codes.Code(err.Code), err.Msg)
}

// ParseGrpcError 解析grpc错误信息
func ParseGrpcError(err error) (common.BusinessCode, string) {
	fromError, _ := status.FromError(err)
	return common.BusinessCode(fromError.Code()), fromError.Message()
}
