/*
创建者：     Zero
创建时间：   2021/5/20
项目名称：   golang-study
*/
package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

//文档地址：https://github.com/gin-contrib/sessions

var sessionName string
var sessionValue string

type MyOption struct {
	sessions.Options
}

func sessionFunc(c *gin.Context) {
	name := c.Query("name")
	if len(name) <= 0 {
		c.JSON(http.StatusBadRequest, "数据错误！")
		return
	}
	sessionName = "session_" + name
	sessionValue = "session_value_" + name
	session := sessions.Default(c)
	sessionData := session.Get(sessionName)
	if sessionData != sessionValue {
		//保存session
		session.Set(sessionName, sessionValue)
		o := MyOption{}
		o.Path = "/"
		o.MaxAge = 10 //有效期，单位s
		session.Options(o.Options)
		session.Save() //保存session
		c.JSON(http.StatusOK, "首次访问，session已保存")
		return
	}
	c.JSON(http.StatusOK, "访问成功，您的session是："+sessionData.(string))
}
