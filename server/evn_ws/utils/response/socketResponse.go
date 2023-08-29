package response

import (
	"dragonsss.cn/evn_common"
	consts "dragonsss.com/evn_ws/utils"
	"dragonsss.com/evn_ws/utils/proto/pb"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

/*json 交互*/

type DataWs struct {
	Code    common.MyCode `json:"code"`
	Type    string        `json:"type,omitempty"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data,omitempty"` // omitempty当data为空时,不展示这个字段
}

func NotLoginWs(ws *websocket.Conn, msg string) {
	rd := &DataWs{
		Code:    common.CodeNotLogin,
		Message: msg,
		Data:    nil,
	}
	err := ws.WriteJSON(rd)
	if err != nil {
		return
	}
}

func SuccessWs(ws *websocket.Conn, tp string, data interface{}) {
	rd := &DataWs{
		Code:    common.CodeSuccess,
		Type:    tp,
		Message: common.CodeSuccess.Msg(),
		Data:    data,
	}
	err := ws.WriteJSON(rd)
	if err != nil {
		return
	}
}

func ErrorWs(ws *websocket.Conn, msg string) {
	rd := &DataWs{
		Code:    common.CodeServerBusy,
		Type:    consts.VideoSocketTypeError,
		Message: msg,
		Data:    nil,
	}
	err := ws.WriteJSON(rd)
	if err != nil {
		return
	}
}

/*proto 交互*/

func ErrorWsProto(ws *websocket.Conn, msg string) {
	message := &pb.Message{
		MsgType: consts.Error,
		Data:    []byte(msg),
	}
	res, _ := proto.Marshal(message)
	_ = ws.WriteMessage(websocket.BinaryMessage, res)
}
