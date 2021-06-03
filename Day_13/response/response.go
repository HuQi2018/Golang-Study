/*
创建者：     Zero
创建时间：   2021/5/25
项目名称：   golang-study
*/
package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(c *gin.Context, httpStatus int, code int, data interface{}, msg string) {
	c.JSON(httpStatus, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

func Success(c *gin.Context, data interface{}, msg string) {
	Response(c, http.StatusOK, 200, data, msg)
}

func Fail(c *gin.Context, data interface{}, msg string) {
	Response(c, http.StatusOK, 400, data, msg)
}
