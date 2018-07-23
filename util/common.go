package util

import (
	"crypto/aes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	ecb "github.com/haowanxing/go-aes-ecb"
	uuid "github.com/satori/go.uuid"
)

type GError struct {
	Code float64
	Err  string
}

func (err *GError) Error() string {
	return fmt.Sprintf("{Code:%d,Err:%s}", err.Code, err.Err)
}

type ResponseMesage struct {
	Status  bool
	Message string
	Error   error
	Data    interface{}
}

func JSON(c *gin.Context, message ResponseMesage) {
	if message.Error != nil {
		message.Status = false
		message.Message = message.Message + "失败"
	} else {
		message.Status = true
		message.Message = message.Message + "成功"
	}

	c.JSON(http.StatusOK, message)
}

func UUID() string {
	return uuid.Must(uuid.NewV4()).String()
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str)) // 需要加密的字符串为
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr) // 输出加密结果
}

func RandNumber(len int) string {
	var result string
	s := rand.NewSource(time.Now().UnixNano())
	rd := rand.New(s)
	for j := 0; j < len; j++ {
		result += strconv.Itoa(rd.Intn(10))
	}

	fmt.Println("aaaaaaa", *(*string)(unsafe.Pointer(&result)))
	return result
}

//生成随机字符串
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
func GetLocation() *time.Location {
	l, _ := time.LoadLocation("Asia/Chongqing")
	return l
}
func AESDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("invalid decrypt key")
	}

	blockSize := block.BlockSize()
	if len(crypted)%blockSize != 0 {
		return nil, errors.New("解密数据格式错误")
	}
	origin, err := ecb.AesDecrypt(crypted, []byte(key)) // ECB解密
	if err != nil {
		return nil, err
	}
	// 使用PKCS#7对解密后的内容去除填充
	origin = ecb.PKCS7UnPadding(origin)
	return origin, nil
}
