package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {

	err := initLogrus() //初始化Logrus日志服务
	if err != nil {
		fmt.Println(err)
		return
	}

	r := gin.Default()

	//日志记录
	r.GET("/logrus1", logrusFunc1)

	//日志记录和日志切割中间件
	r.Use(LogMiddleware())
	r.GET("/logrus2", logrusFunc2)

	//中间件对应的包：github.com/unrolled/secure
	r.Use(HttpsHandler())
	r.GET("/https_get", httpsFunc)

	r.Use(AuthMiddleware())
	r.GET("/auth1", authFunc1)

	//使用https运行
	err = r.RunTLS(":9090", "huqi/Day_11/crt/ca.crt", "huqi/Day_11/crt/ca.key")

	if err != nil {
		fmt.Println("服务开启失败！", err)
	}
}
