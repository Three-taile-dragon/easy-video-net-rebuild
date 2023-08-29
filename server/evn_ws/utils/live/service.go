package live

import (
	"context"
	"dragonsss.cn/evn_common/conversion"
	"dragonsss.com/evn_ws/config"
	"dragonsss.com/evn_ws/internal/dao"
	consts "dragonsss.com/evn_ws/utils"
	"dragonsss.com/evn_ws/utils/proto/pb"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"strconv"
)

// 发送弹幕信息
func serviceSendBarrage(lre LiveRoomEvent, text []byte) error {
	barrageInfo := &pb.WebClientSendBarrageReq{}
	if err := proto.Unmarshal(text, barrageInfo); err != nil {
		return fmt.Errorf("消息格式错误")
	}
	src, _ := conversion.FormattingJsonSrc(lre.Channel.UserInfo.Photo, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	response := &pb.WebClientSendBarrageRes{
		UserId:   float32(lre.Channel.UserInfo.ID),
		Username: lre.Channel.UserInfo.Username,
		Avatar:   src,
		Text:     barrageInfo.Text,
		Color:    barrageInfo.Color,
		Type:     barrageInfo.Type,
	}
	data, err := proto.Marshal(response)
	if err != nil {
		return fmt.Errorf("消息格式错误")
	}
	//将弹幕存入最近消息
	str := conversion.Bytes2String(data)
	rc := dao.Rc
	ctx := context.Background()
	if studentLen, _ := rc.R().LLen(ctx, consts.LiveRoomHistoricalBarrage+strconv.Itoa(int(lre.RoomID))).Result(); studentLen >= 10 {
		err := rc.R().RPop(ctx, consts.LiveRoomHistoricalBarrage+strconv.Itoa(int(lre.RoomID))).Err()
		if err != nil {
			return err
		}
	}
	//消息不足20条 直接插入
	err = rc.R().LPush(ctx, consts.LiveRoomHistoricalBarrage+strconv.Itoa(int(lre.RoomID)), str).Err()
	if err != nil {
		zap.L().Error("房间ID： " + strconv.Itoa(int(lre.RoomID)) + " 最近弹幕存入Redis失败 消息： " + string(data))
		return err
	}
	//格式化响应
	message := &pb.Message{
		MsgType: consts.WebClientBarrageRes,
		Data:    data,
	}
	res, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("消息格式错误")
	}
	for _, v := range Severe.LiveRoom[lre.RoomID] {
		v.MsgList <- res
	}

	return nil
}

// 用户上下线提现
func serviceOnlineAndOfflineRemind(lre LiveRoomEvent, isOnlineOndOffline bool) error {
	//得到当前所有用户
	type userListStruct []*pb.EnterLiveRoom
	userList := make(userListStruct, 0)
	src, _ := conversion.FormattingJsonSrc(lre.Channel.UserInfo.Photo, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	for _, v := range Severe.LiveRoom[lre.RoomID] {
		itemSrc, _ := conversion.FormattingJsonSrc(v.UserInfo.Photo, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
		item := &pb.EnterLiveRoom{
			UserId:   float32(v.UserInfo.ID),
			Username: v.UserInfo.Username,
			Avatar:   itemSrc,
		}
		userList = append(userList, item)
	}

	for i := 0; i < len(Severe.LiveRoom[lre.RoomID]); i++ {
	}
	response := &pb.WebClientEnterLiveRoomRes{
		UserId:   float32(lre.Channel.UserInfo.ID),
		Username: lre.Channel.UserInfo.Username,
		Avatar:   src,
		Type:     isOnlineOndOffline,
		List:     userList,
	}

	//响应输出
	data, err := proto.Marshal(response)
	if err != nil {
		return fmt.Errorf("消息格式错误")
	}
	message := &pb.Message{
		MsgType: consts.WebClientEnterLiveRoomRes,
		Data:    data,
	}
	res, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("消息格式错误")
	}
	for _, v := range Severe.LiveRoom[lre.RoomID] {
		v.MsgList <- res
	}
	return nil
}

// 响应历史消息弹幕
func serviceResponseLiveRoomHistoricalBarrage(lre LiveRoomEvent) error {
	//得到历史消息
	rc := dao.Rc
	ctx := context.Background()
	val, err := rc.R().LRange(ctx, consts.LiveRoomHistoricalBarrage+strconv.Itoa(int(lre.RoomID)), 0, -1).Result()

	if err != nil {
		return fmt.Errorf("获取历史弹幕失败")
	}
	historicalBarrage := &pb.WebClientHistoricalBarrageRes{}
	list := make([]*pb.WebClientSendBarrageRes, 0)
	for _, v := range val {
		barrage := &pb.WebClientSendBarrageRes{}
		if err := proto.Unmarshal(conversion.String2Bytes(v), barrage); err != nil {
			return fmt.Errorf("消息格式错误")
		}
		list = append(list, barrage)
	}
	historicalBarrage.List = list
	data, err := proto.Marshal(historicalBarrage)
	if err != nil {
		return fmt.Errorf("消息格式错误")
	}
	//格式化响应
	message := &pb.Message{
		MsgType: consts.WebClientHistoricalBarrageRes,
		Data:    data,
	}
	res, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("消息格式错误")
	}
	for _, v := range Severe.LiveRoom[lre.RoomID] {
		v.MsgList <- res
	}

	return nil
}
