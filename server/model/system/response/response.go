package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ERROR   = 7
	SUCCESS = 0
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}
