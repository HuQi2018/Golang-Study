/*
创建者：     Zero
创建时间：   2021/5/20
项目名称：   golang-study
*/
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"net/http"
	"unicode/utf8"
)

/*  复杂参数校验
添加依赖：
	go get github.com/go-playground/validator
结构体标签添加限定条件
	类似于：Phone string `validate:"required, numeric, len=11"`  //数字类型，长度为11
验证：
	validate := validator.New()  //初始化（赋值）
	err := validate.Struct(address)
*/

//更多校验类型请阅读：https://github.com/go-playground/validator
type UserInfo struct {
	Id   string `validate:"uuid" json:"id"`           //UUID类型
	Name string `validate:"checkName" json:"name"`    //自定义校验
	Age  uint8  `validate:"min=0,max=130" json:"age"` //0<=Age<=130
}

//含有嵌套类型校验
type ValUser struct {
	Name    string       `validate:"required" json:"name"`        //非空
	Age     uint8        `validate:"gte=0,lte=130" json:"age"`    //0<=age<=130
	Email   string       `vaildate:"required,email" json:"email"` //非空，email格式
	Address []ValAddress `validate:"dive" json:"address"`         //dive关键字代表进入到嵌套结构体进行判断校验
}

type ValAddress struct {
	Province string `json:"province" validate:"required"`    //非空
	City     string `json:"city" validate:"required"`        //非空
	Phone    string `json:"phone" validate:"numeric,len=11"` //numeric数字类型，长度为11
}

//校验注册
var validate *validator.Validate

//初始化验证插件，init与main会自动加载
func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("checkName", checkNameFunc)
}

func validateUser(u ValUser) bool {
	err := validate.Struct(u)
	if err != nil {
		//断言为：validator.ValidationErrors，类型为：[]FieldError
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Println("错误字段：", e.Field())
			fmt.Println("错误的值：", e.Value())
			fmt.Println("错误的tag：", e.Tag())
		}
		return false
	}
	return true
}

func checkNameFunc(fl validator.FieldLevel) bool {
	//字符转换，使得不管输入的是英文还是中文，都得是2-12个
	count := utf8.RuneCountInString(fl.Field().String())
	if count >= 2 && count <= 12 {
		return true
	}
	return false
}

var user UserInfo

func validateFunc1(c *gin.Context) {
	err := c.Bind(&user)
	if err != nil {
		fmt.Println("绑定失败！")
	}

	//使用uuid：go get github.com/satori/go.uuid
	//u1 := uuid.Must(uuid.NewV4(), err)
	//fmt.Println("UUID的值：", u1)
	//user.Id = u1.String()

	//对结构体校验
	err = validate.Struct(user)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Println("错误字段：", e.Field())
			fmt.Println("错误的值：", e.Value())
			fmt.Println("错误的tag：", e.Tag())
		}
		c.JSON(http.StatusBadRequest, "数据校验失败！")
		return
	}
	c.JSON(http.StatusOK, "数据校验成功！")
}

var user2 ValUser

func validateFunc2(c *gin.Context) {
	err := c.Bind(&user2)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, "参数错误，绑定失败！")
		return
	}
	//执行参数的校验
	if validateUser(user2) {
		c.JSON(http.StatusOK, "数据检验成功！")
		return
	}
	c.JSON(http.StatusBadRequest, "数据校验失败！")
}
