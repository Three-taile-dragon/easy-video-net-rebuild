package socket

import (
	"dragonsss.com/evn_ws/utils/chat"
	"dragonsss.com/evn_ws/utils/chatUser"
	"dragonsss.com/evn_ws/utils/live"
	"dragonsss.com/evn_ws/utils/notice"
	sokcet "dragonsss.com/evn_ws/utils/video"
)

func init() {
	//初始化所有socket
	go live.Severe.Start()
	go sokcet.Severe.Start()
	go notice.Severe.Start()
	go chat.Severe.Start()
	go chatUser.Severe.Start()
}
