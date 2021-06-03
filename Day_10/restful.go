/*
创建者：     Zero
创建时间：   2021/5/20
项目名称：   golang-study
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

/**
* 发送POST请求Api接口
* @param url：	请求地址
		data：	POST请求提交的数据
		contentType：	请求体格式，如：application/json
		content：		请求放回的内容
* @return	byte字节流数据切片、error错误信息
* @date 2021/5/20 10:26
*/
func getRestfulAPI(url string, data interface{}, contentType string) ([]byte, error) {
	//创建调用API接口的Client
	client := &http.Client{Timeout: 5 * time.Second}
	//将数据转换为JSON，序列化
	jsonStr, _ := json.Marshal(data)
	//发送请求 注意结构体的请求名字要一致
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("调用API接口出现了错误！")
		return nil, err
	}
	//拿到所有请求到的结果，并返回
	res, err := ioutil.ReadAll(resp.Body)
	return res, err
}

//调用第三方接口的返回数据
type Message struct {
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

//调用第三方接口的请求数据
type UserAPI struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

//测试调用API
func testApi() {
	url := "http://127.0.0.1:9091/login"
	//封装请求的数据
	user := UserAPI{"user", "123456"}
	//调用请求的方法
	data, err := getRestfulAPI(url, user, "application/json")
	fmt.Println(data, err)
	var message Message
	//重新将数据反序列化，注意是地址引用
	err = json.Unmarshal(data, &message)
	if err != nil {
		fmt.Println("数据转换错误！", err)
	}
	fmt.Println(message.Msg, message.Data)
}

//客户端提交请求的数据
type ClientRequest struct {
	UserName string      `json:"user_name"`
	Password string      `json:"password"`
	Other    interface{} `json:"other"`
}

//返回客户端的数据
type ClientResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func getOtherApi(c *gin.Context) {
	var requestData ClientRequest
	var response ClientResponse
	err := c.Bind(&requestData)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.Msg = "请求的参数错误！"
		response.Data = err
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//请求第三方API接口数据
	url := "http://127.0.0.1:9091/login"
	//封装请求的数据
	user := UserAPI{requestData.UserName, requestData.Password}
	//调用请求的方法
	data, err := getRestfulAPI(url, user, "application/json")
	fmt.Println(data, err)
	var message Message
	//重新将数据反序列化，注意是地址引用
	err = json.Unmarshal(data, &message)
	if err != nil {
		fmt.Println("数据转换错误！", err)
	}
	fmt.Println(message.Msg, message.Data)
	response.Code = http.StatusOK
	response.Msg = "请求数据成功"
	response.Data = message.Data
	c.JSON(http.StatusOK, response)

}
