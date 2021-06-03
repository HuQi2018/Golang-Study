/*
创建者：     Zero
创建时间：   2021/5/26
项目名称：   goWeb
*/
package middleware

import (
	"MyGoProject/global"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.Fail(c, nil, fmt.Sprint(err))
				c.Abort()
				return
			}
		}()
	}
}
