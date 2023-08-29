package service

import (
	"context"
	"dragonsss.cn/evn_api/api/other/rpc"
	"dragonsss.cn/evn_api/config"
	other2 "dragonsss.cn/evn_api/pkg/model/other"
	"dragonsss.cn/evn_api/pkg/util/location"
	"dragonsss.cn/evn_common/model/upload"
	"dragonsss.cn/evn_common/validator"
	"dragonsss.cn/evn_grpc/other"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	//Temporary 文件文件存续位置
	Temporary = filepath.ToSlash(config.C.UP.LocalConfig.TmpFileUrl + "assets/tmp")
)

func UploadOss(file *multipart.FileHeader, interface1 string, ctx *gin.Context) (results interface{}, err error) {
	mForm := ctx.Request.MultipartForm
	//上传文件
	var fileName string
	fileName = strings.Join(mForm.Value["name"], fileName)
	var fileInterface string
	fileInterface = strings.Join(mForm.Value["interface"], fileInterface)

	//调用grpc
	//对grpc进行两秒超时处理
	context1, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &other.UploadingDirRequest{
		Interface: fileInterface,
	}
	rsp, err := rpc.OtherServiceClient.UploadingDir(context1, msg)
	if err != nil {
		return nil, fmt.Errorf("上传接口不存在")
	}
	method := new(upload.Upload)
	method.Path = config.C.UP.LocalConfig.TmpFileUrl + rsp.Path
	index := strings.LastIndex(fileName, ".")
	suffix := fileName[index:]
	err = validator.CheckVideoSuffix(suffix)
	if err != nil {
		return nil, fmt.Errorf("非法后缀！")
	}
	//创建临时文件存储
	if !location.IsDir(Temporary) {
		if err = os.MkdirAll(Temporary, 0775); err != nil {
			zap.L().Error("创建临时文件报错路径失败 创建路径为："+method.Path, zap.Error(err))
			return nil, fmt.Errorf("创建保存路径失败")
		}
	}
	tmpDst := filepath.ToSlash(Temporary + "/" + fileName)
	err = ctx.SaveUploadedFile(file, tmpDst)
	if err != nil {
		zap.L().Error("临时文件保存失败-保存路径为："+tmpDst+"错误原因 : ", zap.Error(err))
		return nil, fmt.Errorf("上传失败")
	}
	// 如果interface 是 视频投稿的接口 则只保存到临时文件
	//	后续上传到腾讯云 云点播 使用 createVideoContribution 接口完成
	if interface1 == "videoContribution" || interface1 == "videoContributionCover" {
		return rsp.Path + "/" + fileName, nil
	}

	//stsCredentialResult, err := tengcentyunOss.GetStsCredentialResult()
	bucketUrl, _ := url.Parse(config.C.UP.TencentConfig.Host)
	//bucketUrl, _ := url.Parse("https://easy-video-1300278197.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: bucketUrl}
	c := cos.NewClient(b, &http.Client{
		//设置超时时间
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 COS_SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			//SecretID: stsCredentialResult.Credentials.TmpSecretID,
			//// 环境变量 COS_SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			//SecretKey: stsCredentialResult.Credentials.TmpSecretKey,
			//SessionToken: stsCredentialResult.Credentials.SessionToken,
			SecretID:  config.C.UP.TencentConfig.SecretId,
			SecretKey: config.C.UP.TencentConfig.SecretKey,
			Expire:    time.Duration(time.Hour.Seconds()),
			// Debug 模式，把对应 请求头部、请求内容、响应头部、响应内容 输出到标准输出
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  false,
				RequestBody:    false,
				ResponseHeader: false,
				ResponseBody:   false,
			},
		},
	})

	cTime, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Case2 多线程上传对象，查看上传进度
	opt := &cos.MultiUploadOptions{
		ThreadPoolSize: 5, // 默认五个线程上传
	}
	// 使用临时令牌 sts
	//httpHead := &http.Header{}
	//httpHead.Set("x-cos-security-token", stsCredentialResult.Credentials.SessionToken)

	opt.OptIni = &cos.InitiateMultipartUploadOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			Listener: &cos.DefaultProgressListener{}, // TODO 监听上传情况
			//XOptionHeader: httpHead,
		},
	}
	// c.Object.Upload 不能用 sls 临时令牌
	_, _, upErr := c.Object.Upload(
		cTime, rsp.Path+"/"+fileName, tmpDst, opt,
	)
	//fmt.Println(v)
	if upErr != nil {
		zap.L().Error("保存文件失败-保存路径为："+tmpDst+"错误原因:", zap.Error(upErr))
		return nil, fmt.Errorf("上传失败")
	} else {
		return rsp.Path + "/" + fileName, nil
	}
}

func Upload(file *multipart.FileHeader, ctx *gin.Context) (results interface{}, err error) {
	//如果文件大小超过maxMemory,则使用临时文件来存储multipart/form中文件数据
	err = ctx.Request.ParseMultipartForm(128)
	if err != nil {
		return
	}
	mForm := ctx.Request.MultipartForm
	//上传文件明
	var fileName string
	fileName = strings.Join(mForm.Value["name"], fileName)
	var fileInterface string
	fileInterface = strings.Join(mForm.Value["interface"], fileInterface)

	//调用grpc
	//对grpc进行两秒超时处理
	context1, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &other.UploadingDirRequest{
		Interface: fileInterface,
	}
	rsp, err := rpc.OtherServiceClient.UploadingDir(context1, msg)
	if err != nil {
		return nil, fmt.Errorf("上传接口不存在")
	}
	method := new(upload.Upload)
	method.Path = config.C.UP.LocalConfig.FileUrl + rsp.Path
	index := strings.LastIndex(fileName, ".")
	suffix := fileName[index:]
	err = validator.CheckVideoSuffix(suffix)
	if err != nil {
		return nil, fmt.Errorf("非法后缀！")
	}
	//检测文件夹是否创建
	if !location.IsDir(method.Path) {
		if err = os.MkdirAll(method.Path, 0775); err != nil {
			zap.L().Error("文件夹创建失败 创建路径为："+method.Path, zap.Error(err))
			return nil, fmt.Errorf("文件夹创建失败")
		}
	}
	dst := filepath.ToSlash(method.Path + "/" + fileName)
	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		zap.L().Error("文件保存失败-保存路径为："+dst+"错误原因 : ", zap.Error(err))
		return nil, fmt.Errorf("上传失败")
	} else {
		return rsp.Path + "/" + fileName, nil
	}

}

func UploadSlice(file *multipart.FileHeader, ctx *gin.Context) (results interface{}, err error) {
	//如果文件大小超过maxMemory,则使用临时文件来存储multipart/form中文件数据
	err = ctx.Request.ParseMultipartForm(128)
	if err != nil {
		return
	}
	mForm := ctx.Request.MultipartForm
	//上传文件明
	var fileName string
	fileName = strings.Join(mForm.Value["name"], fileName)
	var fileInterface string
	fileInterface = strings.Join(mForm.Value["interface"], fileInterface)

	//调用grpc
	//对grpc进行两秒超时处理
	context1, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &other.UploadingDirRequest{
		Interface: fileInterface,
	}
	rsp, err := rpc.OtherServiceClient.UploadingDir(context1, msg)
	if err != nil {
		return nil, fmt.Errorf("上传接口不存在")
	}
	method := new(upload.Upload)
	method.Path = config.C.UP.LocalConfig.FileUrl + rsp.Path
	//检测文件夹是否创建
	if !location.IsDir(Temporary) {
		if err = os.MkdirAll(Temporary, 0775); err != nil {
			zap.L().Error("文件夹创建失败 创建路径为："+method.Path, zap.Error(err))
			return nil, fmt.Errorf("文件夹创建失败")
		}
	}
	dst := filepath.ToSlash(Temporary + "/" + fileName)
	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		zap.L().Error("分片上传保存失败-保存路径为："+dst+"错误原因 : ", zap.Error(err))
		return nil, fmt.Errorf("上传失败")
	} else {
		_ = os.Chmod(dst, 0775)
		return rsp.Path + "/" + fileName, nil
	}

}

func UploadCheck(data *other2.UploadCheckStruct) (results interface{}, err error) {
	//调用grpc
	//对grpc进行两秒超时处理
	context1, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &other.UploadingDirRequest{
		Interface: data.Interface,
	}
	rsp, err := rpc.OtherServiceClient.UploadingDir(context1, msg)
	if err != nil {
		return nil, fmt.Errorf("上传接口不存在")
	}
	method := new(upload.Upload)
	method.Path = config.C.UP.LocalConfig.FileUrl + rsp.Path

	list := make(other2.UploadSliceList, 0)
	path := filepath.ToSlash(method.Path + "/" + data.FileMd5)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		//文件已存在
		//zap.L().Info("上传文件 " + data.FileMd5 + " 已存在")
		return location.UploadCheckResponse(true, list, path)
	}
	//取出未上传的分片
	for _, v := range data.SliceList {
		if _, err := os.Stat(filepath.ToSlash(Temporary + "/" + v.Hash)); os.IsNotExist(err) {
			list = append(list, other2.UploadSliceInfo{
				Index: v.Index,
				Hash:  v.Hash,
			})
		}
	}
	return location.UploadCheckResponse(false, list, "")
}

func UploadMerge(data *other2.UploadMergeStruct) (results interface{}, err error) {
	//调用grpc
	//对grpc进行两秒超时处理
	context1, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	msg := &other.UploadingDirRequest{
		Interface: data.Interface,
	}
	rsp, err := rpc.OtherServiceClient.UploadingDir(context1, msg)
	if err != nil {
		return nil, fmt.Errorf("上传接口不存在")
	}
	method := new(upload.Upload)
	method.Path = config.C.UP.LocalConfig.FileUrl + rsp.Path
	if !location.IsDir(method.Path) {
		if err = os.MkdirAll(method.Path, 0775); err != nil {
			zap.L().Error("文件夹创建失败 创建路径为："+method.Path, zap.Error(err))
			return nil, fmt.Errorf("创建保存路径失败")
		}
	}
	dst := filepath.ToSlash(method.Path + "/" + data.FileName)
	list := make(other2.UploadSliceList, 0)
	path := filepath.ToSlash(method.Path + "/" + data.FileName)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		//文件已存在直接返回
		return dst, nil
	}

	//取出未上传的分片
	for _, v := range data.SliceList {
		if _, err := os.Stat(filepath.ToSlash(Temporary + "/" + v.Hash)); os.IsNotExist(err) {
			list = append(list, other2.UploadSliceInfo{
				Index: v.Index,
				Hash:  v.Hash,
			})
		}
	}
	if len(list) > 0 {
		zap.L().Error("上传文件 " + data.FileName + " 分片未全部上传")
		return nil, fmt.Errorf("分片未全部上传")
	}
	cf, err := os.Create(dst)
	if err != nil {
		zap.L().Error("创建的合并后文件失败,err: ", zap.Error(err))
	}
	if err := cf.Close(); err != nil {
		zap.L().Error("创建的合并后文件释放内存失败,err: ", zap.Error(err))
	}
	fileInfo, err := os.OpenFile(dst, os.O_APPEND|os.O_WRONLY, os.ModeSetuid)
	if err != nil {
		zap.L().Error("打开创建的合并后文件失败  path "+dst+" err: ", zap.Error(err))
	}
	defer func(fileInfo *os.File) {
		if err := fileInfo.Close(); err != nil {
			zap.L().Error("关闭资源 err: ", zap.Error(err))
		}
	}(fileInfo)
	//合并操作
	for _, v := range data.SliceList {
		tmpFile, err := os.OpenFile(filepath.ToSlash(Temporary+"/"+v.Hash), os.O_RDONLY, os.ModePerm)
		if err != nil {
			zap.L().Error("合并操作打开临时分片失败 err: ", zap.Error(err))
			break
		}
		b, err := io.ReadAll(tmpFile)
		if err != nil {
			zap.L().Error("合并操作读取分片失败 err: ", zap.Error(err))
			break
		}
		if _, err := fileInfo.Write(b); err != nil {
			zap.L().Error("合并分片追加错误 err: ", zap.Error(err))
			return nil, fmt.Errorf("合并分片追加错误")
		}
		// 关闭分片
		if err := tmpFile.Close(); err != nil {
			zap.L().Error("关闭分片错误 err: ", zap.Error(err))
		}
		if err := os.Remove(tmpFile.Name()); err != nil {
			zap.L().Error("合并操作删除临时分片失败 err: ", zap.Error(err))
		}
	}
	return rsp.Path + "/" + data.FileName, nil
}
