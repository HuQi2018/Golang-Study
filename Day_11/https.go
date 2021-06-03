/*
创建者：     Zero
创建时间：   2021/5/21
项目名称：   golang-study
*/
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"net/http"
)

/*
	1、申请证书
		证书需要收费，不过可以通过KeyManager来生成测试证书（https://keymanager.org/）
	2、核心代码
		router.Use(TlsHandler())  //证书中间件
		router.RunTLS(":8090", path + "jz.crt", path + "jz.key")  //开启TLS（证书服务）的Server
*/

func HttpsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddle := secure.New(secure.Options{
			SSLRedirect: true, //只允许https请求
			//SSLHost: "",  //http到https的重定向
			STSSeconds:           1536000, //Strict-Transport-Security  header的时效：1年
			STSIncludeSubdomains: true,    //includeSubdomains will be appended to the Strict-Transport-Security header
			STSPreload:           true,    //STS Preload（预加载）
			FrameDeny:            true,    //X-Frame-Options 有三个值:DENY（表示该页面不允许在 frame 中展示，即便是在相同域名的页面中嵌套也不允许）、
			// SAMEORIGIN、ALLOW-FROM uri
			ContentTypeNosniff: true, //禁用浏览器的类型猜测行为,防止基于 MIME 类型混淆的攻击
			BrowserXssFilter:   true, //启用XSS保护,并在检查到XSS攻击时，停止渲染页面
			//IsDevelopment: true,   //开发模式
		})
		err := secureMiddle.Process(c.Writer, c.Request)
		//如果不安全，禁止
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "请求方式不安全！")
			return
		}
		//如果是重定向，终止
		if status := c.Writer.Status(); status > 300 && status < 399 {
			c.Abort()
			return
		}
	}
}

type HttpRes struct {
	Code   int    `json:"code"`
	Result string `json:"result"`
}

func httpsFunc(c *gin.Context) {
	fmt.Println(c.Request.Host)
	c.JSON(http.StatusOK, HttpRes{
		Code:   http.StatusOK,
		Result: "请求成功",
	})
}
