package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/augustazz/camellia/logger"
	"io/ioutil"
	"os"
	"strings"
)

var SecretCache = make(map[string][]byte)

type SecretKeyType uint8

const (
	SecretKeyPrv SecretKeyType = iota //private secret
	SecretKeyPub
)

func GetPrvRsaKey(path string) []byte {
	key, err := getRsaKey(secretFileName(path, SecretKeyPrv))
	if err != nil {
		logger.Error(err)
		return nil
	}
	return key
}

func GetPubRsaKey(path string) []byte {
	key, err := getRsaKey(secretFileName(path, SecretKeyPub))
	if err != nil {
		logger.Error(err)
		return nil
	}
	return key
}

func getRsaKey(keyFile string) ([]byte, error) {
	if !FileExists(keyFile) {
		return nil, os.ErrNotExist
	}

	if s, ok := SecretCache[keyFile]; ok {
		return s, nil
	}

	s, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	SecretCache[keyFile] = s
	return s, nil
}

//GenRsaKey 生成RSA公钥私钥
func GenRsaKey(savePath string) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvKey := pem.EncodeToMemory(block)
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubKey := pem.EncodeToMemory(block)

	err = ioutil.WriteFile(secretFileName(savePath, SecretKeyPrv), prvKey, 0644)
	if err != nil {
		logger.Error("private key err: ", err)
		return
	}

	err = ioutil.WriteFile(secretFileName(savePath, SecretKeyPub), pubKey, 0644)
	if err != nil {
		logger.Error("public key err: ", err)
		return
	}
}

//RsaSignWithSha256 签名+加密
func RsaSignWithSha256(data []byte, keyBytes []byte) []byte {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error"))
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		logger.Error("ParsePKCS8PrivateKey err: ", err)
		panic(err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		panic(err)
	}

	return signature
}

//RsaVerySignWithSha256 验证
func RsaVerySignWithSha256(data, signData, keyBytes []byte) bool {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signData)
	if err != nil {
		panic(err)
	}
	return true
}

//RsaEncrypt 公钥加密
func RsaEncrypt(data, keyBytes []byte) []byte {
	//解密pem格式的公钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		panic(err)
	}
	return ciphertext
}

//RsaDecrypt 私钥解密
func RsaDecrypt(ciphertext, keyBytes []byte) []byte {
	//获取私钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error!"))
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// 解密
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		panic(err)
	}
	return data
}

func secretFileName(path string, t SecretKeyType) string {
	if path == "" {
		return path
	}
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	if t == SecretKeyPrv {
		return path + "camellia_rsa"
	}
	return path + "camellia_rsa.pub"
}
