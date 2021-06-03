/*
创建者：     Zero
创建时间：   2021/5/20
项目名称：   golang-study
*/
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//1、添加中间件swaggerFiles使用"github.com/swaggo/gin-swagger/swaggerFiles"
//2、给调用方法添加注解
//3、在main.go文件目录下执行：swag init  生成docs目录
//4、运行gin服务，可访问/swagger/index.html地址管理

//登录信息
type Login struct {
	//UserName string `form:"user" binding:"required"`
	//PassWord string `form:"password" binding:"required,min=6,max=12"`

	UserName string `json:"user_name"` //必须存在且非空
	Password string `json:"password"`  //长度大小6-12
	Remark   string `json:"remark"`
}

// @Tags 登录接口
// @Summary 登录
// @Description Login
// @Accept json
// @Produce json
// @Param username query string true "用户名"
// @Param password query string false "密码"
// @Success 200 {string} json "{"code":200, "data": "{"name": "user_name", "password": "password"}", "msg": "OK"}"
// @Router /login [get]
func login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if username == "Zero" && password == "123456" {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "登录成功！",
			"data": "OK",
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"msg":  "登录失败！",
		"data": "error",
	})
	return
}

// @Tags 注册接口
// @Summary 注册
// @Description Register
// @Accept json
// @Produce json
// @Param user_name formData string true "用户名"
// @Param password formData string true "密码"
// @Success 200 {string} json "{"code":200, "data": "{"name": "user_name", "password": "password"}", "msg": "OK"}"
// @Router /register [post]
func register(c *gin.Context) {
	var login Login
	//swagger 表单数据绑定  formData转换
	err := c.BindQuery(&login)
	//err := c.Bind(&login)
	fmt.Println(login.UserName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "绑定失败，参数错误！",
			"data": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "注册成功！",
		"data": "OK",
	})
	return
}
