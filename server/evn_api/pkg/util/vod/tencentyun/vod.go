package tencentyun

import (
	"dragonsss.cn/evn_api/config"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentyun/vod-go-sdk"
)

func vodUpload() error {
	client := &vod.VodUploadClient{}
	//client.SecretId = "your secretId"
	//client.SecretKey = "your secretKey"
	client.SecretId = config.C.UP.SecretId
	client.SecretKey = config.C.UP.SecretKey

	req := vod.NewVodUploadRequest()
	req.MediaFilePath = common.StringPtr("F:/Vue3/easy-video-net-rebuild/server/assets/static/video/users/videoContribution/6ea440907ab908cc801b9253fdb691c94e62bd5700c238568db6b21fd5530b45.mp4")
	//req.StorageRegion = common.StringPtr("ap-guangzhou")

	rsp, err := client.Upload("ap-guangzhou", req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(*rsp.Response.FileId)
	fmt.Println(*rsp.Response.MediaUrl)
	return nil
}
