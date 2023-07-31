package conversion

import (
	"encoding/json"
	"fmt"
)

type Img struct {
	Src string `json:"src" gorm:"column:src"`
	Tp  string `json:"type" gorm:"column:type"`
}

func FormattingJsonSrc(str []byte, localhost string, tencentOssHost string) (url string, err error) {
	data := new(Img)
	err = json.Unmarshal(str, data)
	if err != nil {
		return "", fmt.Errorf("json format error")
	}
	if data.Src == "" {
		return "", nil
	}
	path, err := SwitchIngStorageFun(data.Tp, data.Src, localhost, tencentOssHost)
	if err != nil {
		return "", err
	}
	return path, nil
}

// SwitchIngStorageFun 根据类型拼接路径
func SwitchIngStorageFun(tp string, path string, localhost string, tencentOssHost string) (url string, err error) {
	prefix, err := SwitchTypeAsUrlPrefix(tp, localhost, tencentOssHost)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", prefix, path), nil
}

//TODO 优化配置读取

// SwitchTypeAsUrlPrefix 取url前缀
func SwitchTypeAsUrlPrefix(tp string, localhost string, tencentOssHost string) (url string, err error) {
	switch tp {
	case "local":
		return localhost, nil
	case "tencentOss":
		return tencentOssHost, nil
	default:
		return "", fmt.Errorf("undefined format")
	}
}
