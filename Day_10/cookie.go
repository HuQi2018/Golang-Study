/*
创建者：     Zero
创建时间：   2021/5/20
项目名称：   golang-study
*/
package main

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

var cookieName string
var cookieValue string

//定义cookie中间件
func CookieAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, _ := c.Cookie(cookieName)
		if val == "" {
			c.SetCookie(cookieName, cookieValue, 3600, "/", "localhost", true, true)
		}
	}
}

func cookieFunc1(c *gin.Context) {
	name := c.Query("name")
	if len(name) <= 0 {
		c.JSON(http.StatusBadRequest, "数据错误！")
		return
	}
	cookieName = "cookie_" + name                                  //cookie的key值
	cookieValue = hex.EncodeToString([]byte(cookieName + "value")) //cookie的value值
	val, _ := c.Cookie(cookieName)
	if val == "" {
		c.String(http.StatusOK, "Cookie：%s存储成功！", cookieName)
		return
	}
	c.String(http.StatusOK, "验证成功，cookie值为：%s", val)
}
