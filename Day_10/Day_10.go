package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "golang-study/huqi/Day_10/docs"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

/*
1、调用restful接口
2、参数校验
3、Swagger
4、Cookie
5、Session
*/

var g errgroup.Group

func main() {

	//部署多台服务程序同时运行
	//服务器1：http://127.0.0.1:9091/MulServer01
	server01 := &http.Server{
		Addr:         ":9091",
		Handler:      router01(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	//服务器2：http://127.0.0.1:9092/MulServer02
	//server02 := &http.Server{
	//	Addr:         ":9092",
	//	Handler:      router02(),
	//	ReadTimeout:  5 * time.Second,
	//	WriteTimeout: 10 * time.Second,
	//}
	//开启服务
	g.Go(func() error { //开启服务程序1
		return server01.ListenAndServe()
	})
	//g.Go(func() error {
	//	return server02.ListenAndServe()
	//})
	if err := g.Wait(); err != nil {
		fmt.Println("执行失败：", err)
	}

}

func router01() http.Handler {

	r := gin.Default() //获取路有引擎

	//swagger 中间件主要的作用是：方便前端对接口进行调试，不影响接口的实际功能
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) //使用swagger中间件

	r.GET("/login", login)

	//绑定请求参数与结构体一致
	r.POST("/register", register)
	//testApi()
	//restful接口调用
	r.POST("/getApiData", getOtherApi)
	//简单参数校验
	r.POST("/validate", validateFunc1)
	//复杂参数校验
	r.POST("/validate2", validateFunc2)

	//cookie中间件
	r.Use(CookieAuth())
	r.GET("/cookie", cookieFunc1)

	//session中间件
	store := cookie.NewStore([]byte("session_secret")) //session_secret为服务器私有内容
	r.Use(sessions.Sessions("mysession", store))
	r.GET("/session", sessionFunc)

	return r
}

func router02() http.Handler {
	r := gin.Default()

	return r
}
