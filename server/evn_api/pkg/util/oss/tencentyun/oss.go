package tencentyun

import (
	"context"
	"dragonsss.cn/evn_api/config"
	"go.uber.org/zap"

	//"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"net/http"
	"net/url"
	"time"
)

var SecretId = config.C.UP.TencentConfig.SecretId
var SecretKey = config.C.UP.TencentConfig.SecretKey
var appid = config.C.UP.TencentConfig.Appid
var bucketName = config.C.UP.TencentConfig.Bucket
var region = config.C.UP.TencentConfig.Region
var durationSeconds = config.C.UP.TencentConfig.DurationSeconds
var bucketUrl = config.C.UP.TencentConfig.Host

// STS

func CreateStsClient(secretId, secretKey string) *sts.Client {
	c := sts.NewClient(
		// 通过环境变量获取密钥, os.Getenv 方法表示获取环境变量
		secretId,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考https://cloud.tencent.com/document/product/598/37140
		secretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考https://cloud.tencent.com/document/product/598/37140
		nil,
		// sts.Host("sts.internal.tencentcloudapi.com"), // 设置域名, 默认域名sts.tencentcloudapi.com
		// sts.Scheme("http"),      // 设置协议, 默认为https，公有云sts获取临时密钥不允许走http，特殊场景才需要设置http
	)
	return c
}

func CreateStsCredentialOptions(region string) *sts.CredentialOptions {
	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						"name/cos:PutObject",
						"name/cos:GetService",
						"name/cos:GetBucket",
						"name/cos:HeadService",
						"name/cos:PostObject",
						"name/cos:HeadObject",
						"name/cos:GetObject",
						"name/cos:InitiateMultipartUpload",
						"name/cos:UploadPart",
						"name/cos:UploadPartCopy",
						"name/cos:CompleteMultipartUpload",
						"name/cos:AbortMultipartUpload",
						"name/cos:ListMultipartUploads",
						"name/cos:ListParts",
						"name/cos:CreateJob",
						"name/cos:DescribeJob",
						"name/cos:ListJobs",
					},
					Effect: "allow",
					Resource: []string{
						//这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
						//存储桶的命名格式为 BucketName-APPID，此处填写的 bucket 必须为此格式
						"qcs::cos:" + region + ":uid/" + appid + ":" + bucketName + "/assets/*",
					},
					// 开始构建生效条件 condition
					// 关于 condition 的详细设置规则和COS支持的condition类型可以参考https://cloud.tencent.com/document/product/436/71306
					Condition: map[string]map[string]interface{}{
						//TODO ip限制
						//"ip_equal": map[string]interface{}{
						//	"qcs:ip": []string{
						//		"10.217.182.3/24",
						//		"111.21.33.72/24",
						//	},
						//},
					},
				},
			},
		},
	}
	return opt
}

func GetStsCredentialResult() (*sts.CredentialResult, error) {
	client := CreateStsClient(SecretId, SecretKey)
	option := CreateStsCredentialOptions(region)
	// 策略概述 https://cloud.tencent.com/document/product/436/18023
	res, err := client.GetCredential(option)
	//fmt.Printf("%+v\n", res)
	//fmt.Printf("%+v\n", res.Credentials)
	return res, err
}

// RegisterMediaInfo 注册媒体信息	用于视频转码 可使用腾讯云oss提供的工作流自动处理
//这里先不做处理

// GetMediaInfo 获取媒体信息

func getClient() (*cos.Client, string, error) {
	bucketUrl, err := url.Parse(bucketUrl)
	if err != nil {
		zap.L().Error("获取bucketUrl失败 err", zap.Error(err))
	}
	//ciUrl, _ := url.Parse("https://test-1234567890.ci.ap-chongqing.myqcloud.com")
	stsCredential, _ := GetStsCredentialResult()
	b := &cos.BaseURL{BucketURL: bucketUrl}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			//SecretID:  stsCredential.Credentials.TmpSecretID,
			//SecretKey: stsCredential.Credentials.TmpSecretKey,
			//SessionToken: stsCredential.Credentials.SessionToken,
			SecretID:  SecretId,
			SecretKey: SecretKey,
			//Expire:       time.Duration(durationSeconds),
			// Debug 模式，把对应 请求头部、请求内容、响应头部、响应内容 输出到标准输出
			Transport: &debug.DebugRequestTransport{
				RequestHeader: false,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    false,
				ResponseHeader: false,
				ResponseBody:   false,
			},
		},
	})
	return c, stsCredential.Credentials.SessionToken, err
}

func GetMediaInfo(mediaPath string) (*cos.MediaInfo, error) {
	c, _, err := getClient()
	if err != nil {
		return nil, err
	}
	// 使用 临时令牌 sts
	//httpHead := &http.Header{}
	//httpHead.Set("x-cos-security-token", stsToken)
	//
	//opt := &cos.ObjectGetOptions{
	//	XOptionHeader: httpHead,
	//}
	res, _, err := c.CI.GetMediaInfo(context.Background(), mediaPath, nil)
	if err != nil {
		zap.L().Error("获取媒体信息失败 err", zap.Error(err))
	}
	//fmt.Printf("res: %+v\n", res.MediaInfo)
	return res.MediaInfo, err
}

// SubmitTranscodeJob 视频云转码 可使用腾讯云oss提供的工作流自动处理
//这里先不做处理
