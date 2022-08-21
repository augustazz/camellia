package util

import (
	"encoding/hex"
	"fmt"
	"github.com/augustazz/camellia/config"
	"testing"
)

var connConf = config.GetConnConfig()

func TestGenKey(t *testing.T) {
	fmt.Println("-------------------------------生成RSA公私钥-----------------------------------------")
	GenRsaKey(connConf.AuthFilePath)
}

func TestRsaDecrypt(t *testing.T) {
	//rsa 密钥文件产生
	fmt.Println("-------------------------------获取RSA公私钥-----------------------------------------")
	prvKey := GetPrvRsaKey(connConf.AuthFilePath)
	if prvKey == nil {
		fmt.Println("read err")
		return
	}
	pubKey := GetPubRsaKey("")
	if pubKey == nil {
		fmt.Println("read err")
		return
	}
	fmt.Println(string(prvKey))
	fmt.Println(string(pubKey))

	fmt.Println("-------------------------------进行签名与验证操作-----------------------------------------")
	var data = "100023DT39485XVlBzgbaiCMRAjWwhTHctcuAxhxKQFDaFpLSjFbcXoEFfRsWxPLDnJObCsNVlgTe"
	fmt.Println("对消息进行签名操作...")
	signData := RsaSignWithSha256([]byte(data), prvKey)
	fmt.Println("消息的签名信息： ", hex.EncodeToString(signData))
	fmt.Println("\n对签名信息进行验证...")
	if RsaVerySignWithSha256([]byte(data), signData, pubKey) {
		fmt.Println("签名信息验证成功，确定是正确私钥签名！！")
	}

	fmt.Println("-------------------------------进行加密解密操作-----------------------------------------")
	ciphertext := RsaEncrypt([]byte(data), pubKey)
	fmt.Println("公钥加密后的数据：", hex.EncodeToString(ciphertext))
	sourceData := RsaDecrypt(ciphertext, prvKey)
	fmt.Println("私钥解密后的数据：", string(sourceData))
}
