package common

import (
	"dragonsss.cn/evn_api/pkg/model/other"
	ws2 "dragonsss.cn/evn_grpc/ws"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"net/http"
)

// 定义消息返回体

type BusinessCode int
type Result struct {
	Code    BusinessCode `json:"code"`
	Message string       `json:"message"`
	Data    any          `json:"data"`
}

func (r *Result) Success(data any) *Result {
	r.Code = http.StatusOK
	r.Message = "success"
	r.Data = data
	return r
}

func (r *Result) Fail(code BusinessCode, msg string) *Result {
	r.Code = code
	r.Message = msg
	return r
}

/*json 交互*/

type DataWs struct {
	Code    other.MyCode `json:"code"`
	Type    string       `json:"type,omitempty"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data,omitempty"` // omitempty当data为空时,不展示这个字段
}

func NotLoginWs(ws *websocket.Conn, msg string) {
	rd := &DataWs{
		Code:    other.CodeNotLogin,
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
		Code:    other.CodeSuccess,
		Type:    tp,
		Message: other.CodeSuccess.Msg(),
		Data:    data,
	}
	err := ws.WriteJSON(rd)
	if err != nil {
		return
	}
}

func ErrorWs(ws *websocket.Conn, msg string) {
	rd := &DataWs{
		Code:    other.CodeServerBusy,
		Type:    other.VideoSocketTypeError,
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
	message := &ws2.Message{
		MsgType: other.Error,
		Data:    []byte(msg),
	}
	res, _ := proto.Marshal(message)
	_ = ws.WriteMessage(websocket.BinaryMessage, res)
}
