package model

import (
	"dragonsss.cn/evn_common/conversion"
	"dragonsss.cn/evn_common/model/user"
)

type ReqGetRoom struct {
	Status int    `json:"status" binding:"required"`
	Data   string `json:"data" binding:"required"`
}

type GetLiveRoomInfoReceiveStruct struct {
	RoomID uint `json:"room_id"`
}

type LivestatRes struct {
	Status int `json:"status"`
	Data   struct {
		Publishers []struct {
			Key             string `json:"key"`
			Url             string `json:"url"`
			StreamId        int    `json:"stream_id"`
			VideoTotalBytes int64  `json:"video_total_bytes"`
			VideoSpeed      int    `json:"video_speed"`
			AudioTotalBytes int    `json:"audio_total_bytes"`
			AudioSpeed      int    `json:"audio_speed"`
		} `json:"publishers"`
		Players interface{} `json:"players"`
	} `json:"data"`
}

type GetLiveRoomResponseStruct struct {
	Address string `json:"address"`
	Key     string `json:"key"`
}

func GetLiveRoomResponse(address string, key string) interface{} {
	return GetLiveRoomResponseStruct{
		Address: address,
		Key:     key,
	}
}

type GetLiveRoomInfoResponseStruct struct {
	Username   string `json:"username"`
	Photo      string `json:"photo"`
	LiveTitle  string `json:"live_title"`
	Flv        string `json:"flv"`
	LicenseUrl string `json:"licenseUrl"`
}

func GetLiveRoomInfoResponse(info *user.User, flv string, licenseUrl string, localhost string, tencentOssHost string) interface{} {
	photo, _ := conversion.FormattingJsonSrc(info.Photo, localhost, tencentOssHost)
	return GetLiveRoomInfoResponseStruct{
		Username:   info.Username,
		Photo:      photo,
		LiveTitle:  info.LiveInfo.Title,
		Flv:        flv,
		LicenseUrl: licenseUrl,
	}
}

type BeLiveInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Photo    string `json:"photo"`
	Img      string `json:"img"`
	Title    string `json:"title"`
	Online   int    `json:"online"`
}

type BeLiveInfoList []BeLiveInfo

func GetBeLiveListResponse(ul *user.UserList, localhost string, tencentOssHost string) interface{} {
	list := make(BeLiveInfoList, 0)
	for _, v := range *ul {
		photo, _ := conversion.FormattingJsonSrc(v.Photo, localhost, tencentOssHost)
		img, _ := conversion.FormattingJsonSrc(v.LiveInfo.Img, localhost, tencentOssHost)
		list = append(list, BeLiveInfo{
			ID:       v.ID,
			Username: v.Username,
			Photo:    photo,
			Img:      img,
			Title:    v.LiveInfo.Title,
			Online:   0,
		})
	}
	return list
}
