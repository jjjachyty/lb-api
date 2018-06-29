package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
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
	fmt.Println(cipherStr)
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

func GetLocation() *time.Location {
	l, _ := time.LoadLocation("Asia/Chongqing")
	return l
}
