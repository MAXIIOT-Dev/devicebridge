package controllers

import (
	"github.com/gin-gonic/gin"
)

// ResponseData vbase response
type ResponseData struct {
	ErrorCode int         `json:"errorCode"` // 错误代码
	Status    int         `json:"status"`    // 状态:。0 : 响应成功，1：响应失败
	Message   string      `json:"message"`   // 错误信息
	Result    interface{} `json:"result"`    // 数据
}

// Response response vbase data format
func Response(c *gin.Context, httpStatus int, errorCode int, status int, msg string, result interface{}) {
	c.JSON(httpStatus, ResponseData{
		ErrorCode: errorCode,
		Status:    status,
		Message:   msg,
		Result:    result,
	})
}
