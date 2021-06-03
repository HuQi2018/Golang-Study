/*
创建者：     Zero
创建时间：   2021/5/25
项目名称：   golang-study
*/
package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-study/huqi/Day_13/response"
)

func RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response.Fail(c, nil, fmt.Sprint(err))
				c.Abort()
				return
			}
		}()
	}
}
