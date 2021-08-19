package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
)

// MD5V @author: [piexlmax](https://github.com/piexlmax)
//@function: MD5V
//@description: md5加密
//@param: str []byte
//@return: string
// 知识补充：md5加密是不可逆的，用于写入数据库的加密算法
func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

// RSA非对称加密 解决前端传密码明文问题

// GenerateRSAKey 生成RSA私钥和公钥
func GenerateRSAKey() (err error, X509PrivateKey []byte, X509PublicKey []byte) {
	bits := 2048
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err, nil, nil
	}
	//保存私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey = x509.MarshalPKCS1PrivateKey(privateKey)
	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err = x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err, nil, nil
	}
	return nil, X509PrivateKey, X509PublicKey
}

// RSAEncrypt RSA加密
func RSAEncrypt(plainText []byte, key []byte) []byte {
	publicKeyInterface, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		panic(err)
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)

	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		panic(err)
	}
	//返回密文
	return cipherText
}

// RSADecrypt RSA解密
func RSADecrypt(cipherText string, key []byte) string {
	/* 对前端的密码做处理 这一步非常重要!!!!! */
	cipherTextByte, _ := base64.StdEncoding.DecodeString(cipherText)
	privateKey, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		panic(err)
	}
	//对密文进行解密
	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherTextByte)
	//返回明文
	return string(plainText)
}
