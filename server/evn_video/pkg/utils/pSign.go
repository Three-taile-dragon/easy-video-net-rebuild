package utils

import (
	"dragonsss.com/evn_video/config"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

func PSignCalculate(fileID string) (string, error) {
	appId := config.C.Vod.Appid                                 // 用户 appid
	fileId := fileID                                            // 目标 FileId
	audioVideoType := config.C.Vod.AudioVideoType               // 播放的音视频类型 未加密的 转自适应码流 输出
	rawAdaptiveDefinition := config.C.Vod.RawAdaptiveDefinition // 允许输出的未加密的自适应码流模板 ID
	//imageSpriteDefinition := 10     // 做进度条预览的雪碧图模板 ID
	currentTime := time.Now().Unix()
	psignExpire := currentTime + config.C.Vod.PsignExpire // 可任意设置过期时间，示例1h
	urlTimeExpire := strconv.FormatInt(psignExpire, 16)   // 可任意设置过期时间，16进制字符串形式，示例1h
	playKey := []byte(config.C.Vod.Key)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"appId":  appId,
		"fileId": fileId,
		"contentInfo": map[string]any{
			"audioVideoType":        audioVideoType,
			"rawAdaptiveDefinition": rawAdaptiveDefinition,
		},
		"currentTimeStamp": currentTime,
		"expireTimeStamp":  psignExpire,
		"urlAccessInfo": map[string]string{
			"t": urlTimeExpire,
		},
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(playKey)
	return tokenString, err
}
