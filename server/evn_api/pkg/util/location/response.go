package location

import (
	"dragonsss.cn/evn_api/config"
	other2 "dragonsss.cn/evn_api/pkg/model/other"
	"dragonsss.cn/evn_common/calculate"
	"dragonsss.cn/evn_common/conversion"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/video"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"time"
)

func UploadingMethodResponse(tp string) interface{} {
	type UploadingMethodResponseStruct struct {
		Tp string `json:"type"`
	}
	return UploadingMethodResponseStruct{
		Tp: tp,
	}
}

func UploadingDirResponse(dir string, quality float64) interface{} {
	type UploadingDirResponseStruct struct {
		Path    string  `json:"path"`
		Quality float64 `json:"quality"`
	}
	return UploadingDirResponseStruct{
		Path:    dir,
		Quality: quality,
	}
}

// VideoInfo 首页视频
type VideoInfo struct {
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

type videoInfoList []VideoInfo

func SearchVideoResponse(videoList *video.VideosContributionList, localhost string, tencentOssHost string) (interface{}, error) {
	//处理视频
	vl := make(videoInfoList, 0)
	for _, lk := range *videoList {
		cover, _ := conversion.FormattingJsonSrc(lk.Cover, localhost, tencentOssHost)
		videoSrc, _ := conversion.FormattingJsonSrc(lk.Video, localhost, tencentOssHost)
		info := VideoInfo{
			ID:            lk.ID,
			Uid:           lk.Uid,
			Title:         lk.Title,
			Video:         videoSrc,
			Cover:         cover,
			VideoDuration: lk.VideoDuration,
			Label:         conversion.StringConversionMap(lk.Label),
			Introduce:     lk.Introduce,
			Heat:          lk.Heat,
			BarrageNumber: len(lk.Barrage),
			Username:      lk.UserInfo.Username,
			CreatedAt:     lk.CreatedAt,
		}
		vl = append(vl, info)
	}
	return vl, nil
}

type UserInfo struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Photo       string `json:"photo"`
	Signature   string `json:"signature"`
	IsAttention bool   `json:"is_attention"`
}

type UserInfoList []UserInfo

func SearchUserResponse(userList *user.UserList, aids []uint, localhost string, tencentOssHost string) (interface{}, error) {
	list := make(UserInfoList, 0)
	for _, v := range *userList {
		photo, _ := conversion.FormattingJsonSrc(v.Photo, localhost, tencentOssHost)
		list = append(list, UserInfo{
			ID:          v.ID,
			Username:    v.Username,
			Photo:       photo,
			Signature:   v.Signature,
			IsAttention: calculate.ArrayIsContain(aids, v.ID),
		})
	}
	return list, nil
}

type UploadCheckStruct struct {
	IsUpload bool                   `json:"is_upload"`
	List     other2.UploadSliceList `json:"list"`
	Path     string                 `json:"path"`
}

func UploadCheckResponse(is bool, list other2.UploadSliceList, path string) (interface{}, error) {
	return UploadCheckStruct{
		IsUpload: is,
		List:     list,
		Path:     path,
	}, nil
}

type GteStsInfoStruct struct {
	Region          string `json:"region"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	StsToken        string `json:"sts_token"`
	Bucket          string `json:"bucket"`
	ExpirationTime  int64  `json:"expiration_time"`
}

func GteStsInfo(info *sts.Credentials) (interface{}, error) {
	return GteStsInfoStruct{
		Region:          config.C.UP.TencentConfig.Region,
		AccessKeyID:     info.TmpSecretID,
		AccessKeySecret: info.TmpSecretKey,
		StsToken:        info.SessionToken,
		Bucket:          config.C.UP.TencentConfig.Bucket,
		ExpirationTime:  time.Now().Unix() + int64(config.C.UP.TencentConfig.DurationSeconds),
	}, nil
}
