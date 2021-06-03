/*
创建者：     Zero
创建时间：   2021/5/25
项目名称：   golang-study
*/
package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-study/huqi/Day_13/common"
	"golang-study/huqi/Day_13/model"
	"golang-study/huqi/Day_13/response"
	"golang-study/huqi/Day_13/util"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Register(c *gin.Context) {
	var reqUser model.UserBase
	err := c.Bind(&reqUser)
	if err != nil {
		response.Fail(c, "Error", "请求参数有误！")
		return
	}
	name := reqUser.Name
	telephone := reqUser.Telephone
	password := reqUser.Password
	//数据验证
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位！")
		fmt.Println(telephone, len(telephone))
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位！")
		return
	}
	//如果用户没有填写用户名则系统自动生成随机用户名
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	//判断用户手机号是否已经存在
	var user model.UserBase
	common.MyDB.Where("telephone=?", telephone).First(&user)
	fmt.Println(user.ID)
	if user.ID != 0 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户手机号已经存在，请更换！")
		return
	}
	//创建用户
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "用户密码加密错误！")
		return
	}
	newUser := model.UserBase{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashPassword),
	}
	common.MyDB.Create(&newUser) //新增记录
	//分发密钥
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, err, "系统异常，分发密钥失败！")
		return
	}
	//返回结果
	response.Success(c, gin.H{"tooken": token}, "注册成功！")
}

func Login(c *gin.Context) {
	var reqUser model.UserBase
	err := c.Bind(&reqUser)
	if err != nil {
		response.Fail(c, "Error", "请求参数有误！")
		return
	}
	//name := reqUser.Name
	telephone := reqUser.Telephone
	password := reqUser.Password
	//数据验证
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位！")
		fmt.Println(telephone, len(telephone))
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位！")
		return
	}
	//依据手机号，查询用户注册的数据记录
	var user model.UserBase
	common.MyDB.Where("telephone=?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户不存在！")
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户不存在！")
		return
	}
	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, err, "系统异常，分发密钥失败！")
		return
	}
	//返回结果
	response.Success(c, gin.H{"tooken": token}, "登录成功！")
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")
	response.Success(c, gin.H{
		"user": response.ToUserDto(user.(model.UserBase)),
	}, "获取成功！")
}
