package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()       //路有引擎
	r.GET("/Get", getMsg)    //定义get访问地址
	r.POST("/Post", postMsg) //定义post访问地址
	r.GET("/test1", test1)

	//重定向到其他网址
	r.GET("/RedictToBaiDu", func(c *gin.Context) {
		url := "http://www.baidu.com"
		c.Redirect(http.StatusMovedPermanently, url)
	})
	//内部重定向
	r.GET("/RedictToLocal", func(c *gin.Context) {
		c.Request.URL.Path = "/test1?name=Redict"
		r.HandleContext(c)
	})

	r.Run("127.0.0.1:9090") //如果不指定IP地址、端口号，默认为本地ip地址，8080端口，可直接简写端口号“:9090”
}

func test1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "返回信息",
		"data": "欢迎您！",
	})
}

//Post方式处理的请求
func postMsg(c *gin.Context) {
	name1 := c.DefaultPostForm("name", "Gin")
	fmt.Println(name1)
	name2, _ := c.GetPostForm("name") //返回两个值
	fmt.Println(name2)
	c.String(http.StatusOK, "欢迎您的访问：%s", name1)
}

//Get方式处理的请求
func getMsg(c *gin.Context) {
	name := c.Query("name") //获取URL中的值
	//返回字符串类型
	//c.String(http.StatusOK, "欢迎您的访问：%s", name)
	//返回JSON数据类型
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "返回信息",
		"data": "欢迎您：" + name,
	})
}
