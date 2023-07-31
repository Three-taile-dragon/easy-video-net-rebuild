package common

import (
	common "dragonsss.cn/lbk_common"
	"fmt"
	"gitee.com/frankyu365/gocrypto/ecc"
)

func EccInit() {
	//生成ECC加密密钥
	pem, _ := common.PathExists("F:/Vue3/lbk_background/lbk_common/encrypts/ecc/key/eccPublic.pem") //判断是否存在EccPem
	key, _ := common.PathExists("F:/Vue3/lbk_background/lbk_common/encrypts/ecc/key/eccPublic.key") //判断是否存在EccKey
	if !pem && !key {
		GenerateECCKey() //如果都不存在就重新生成Ecc密钥
	}
}

func GenerateECCKey() {
	err := ecc.GenerateECCKey(256, "F:/Vue3/lbk_background/lbk_common/encrypts/ecc/key/")
	if err != nil {
		fmt.Println(err)
	}
}

func EccEncrypt(plain string) []byte {
	plainText := []byte(plain)
	cipherText, err := ecc.EccEncrypt(plainText, "./lbk_common/encrypts/ecc/key/eccPublic.pem")
	if err != nil {
		fmt.Println(err)
	}
	return cipherText

}

func EccDecrypt(cipher []byte) []byte {
	plainText, err := ecc.EccDecrypt(cipher, "./lbk_common/encrypts/ecc/key/eccPrivate.pem")
	if err != nil {
		fmt.Println(err)
	}
	return plainText
}
