package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"response"`
	Msg  string      `json:"message"`
}

const (
	ERROR       = 7
	SUCCESS     = 0
	EXAMSUCCESS = 1
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func OkWithDetailedExam(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func NoAuth(message string, c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		7,
		nil,
		message,
	})
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithMessageExam(message string, c *gin.Context) {
	Result(EXAMSUCCESS, message, "成功", c)
}
