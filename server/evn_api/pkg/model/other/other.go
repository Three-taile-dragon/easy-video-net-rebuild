package other

import "time"

type HomeRequest struct {
	Page int64 `json:"page"`
	Size int64 `json:"size"`
}

// 首页轮播图
type rotographInfo struct {
	Title string `json:"title"`
	Cover string `json:"cover"`
	Color string `json:"color"`
	Type  string `json:"type"`
	ToId  uint   `json:"to_id"`
}
type RotographInfoList []rotographInfo

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

type VideoInfoList []VideoInfo

type GetLiveRoomInfoReceiveStruct struct {
	RoomID uint `json:"room_id"`
}

type GetLiveRoomInfoResponseStruct struct {
	Username  string `json:"username"`
	Photo     string `json:"photo"`
	LiveTitle string `json:"live_title"`
	Flv       string `json:"flv"`
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
