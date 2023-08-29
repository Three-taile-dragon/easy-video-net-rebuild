package ws

import (
	consts "dragonsss.com/evn_ws/utils"
	"dragonsss.com/evn_ws/utils/chat"
	"dragonsss.com/evn_ws/utils/chatUser"
	"dragonsss.com/evn_ws/utils/live"
	"dragonsss.com/evn_ws/utils/notice"
	"dragonsss.com/evn_ws/utils/proto/pb"
	"dragonsss.com/evn_ws/utils/response"
	sokcet "dragonsss.com/evn_ws/utils/video"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"strconv"
)

type HandleWs struct {
}

func New() *HandleWs {
	return &HandleWs{}
}

func (w *HandleWs) NoticeSocket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	ws := conn.(*websocket.Conn)
	err := notice.CreateNoticeSocket(uid, ws)
	if err != nil {
		response.ErrorWs(ws, "创建通知socket失败")
	}
}

func (w *HandleWs) ChatSocket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	ws := conn.(*websocket.Conn)
	err := chat.CreateChatSocket(uid, ws)
	if err != nil {
		response.ErrorWs(ws, "创建聊天socket失败")
	}
}

func (w *HandleWs) ChatByUserSocket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	//判断是否创建视频socket房间
	tidQuery, _ := strconv.Atoi(ctx.Query("tid"))
	tid := uint(tidQuery)
	ws := conn.(*websocket.Conn)
	err := chatUser.CreateChatByUserSocket(uid, tid, ws)
	if err != nil {
		response.ErrorWs(ws, "创建用户聊天socket失败")
	}
}

func (w *HandleWs) LiveSocket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	ws := conn.(*websocket.Conn)

	//判断是否创建直播间
	liveRoom, _ := strconv.Atoi(ctx.Query("liveRoom"))
	liveRoomID := uint(liveRoom)
	if live.Severe.LiveRoom[liveRoomID] == nil {
		message := &pb.Message{
			MsgType: consts.Error,
			Data:    []byte("直播未开启"),
		}
		res, _ := proto.Marshal(message)
		_ = ws.WriteMessage(websocket.BinaryMessage, res)
		return
	}

	err := live.CreateSocket(ctx, uid, liveRoomID, ws)
	if err != nil {
		response.ErrorWs(ws, err.Error())
		return
	}
}

func (w *HandleWs) VideoSocket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	ws := conn.(*websocket.Conn)
	//判断是否创建视频socket房间
	id, _ := strconv.Atoi(ctx.Query("videoID"))
	videoID := uint(id)
	if sokcet.Severe.VideoRoom[videoID] == nil {
		//无人观看主动创建
		sokcet.Severe.VideoRoom[videoID] = make(sokcet.UserMapChannel, 10)
	}
	err := sokcet.CreateVideoSocket(uid, videoID, ws)
	if err != nil {
		response.ErrorWs(ws, "创建socket失败")
	}
}
