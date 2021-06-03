# Go语言学习笔记 Day08

### [Day07](http://njgit.jsaepay.com/wangwei/golang-study/src/branch/master/huqi/Day_07)
### [Day09](http://njgit.jsaepay.com/wangwei/golang-study/src/branch/master/huqi/Day_09)

------------

####[Go Web 开发](https://www.w3cschool.cn/yqbmht)
####[51CTO学院 Gin框架开发基础与提升](https://edu.51cto.com/center/course/lesson/index?id=672385)

### 1. Beedb管理数据库
	beedb是支持database/sql标准接口的ORM库
		Mysql:github.com/ziutek/mymysql/godrv[*]
		Mysql:code.google.com/p/go-mysql-driver[*]
		PostgreSQL:github.com/bmizerany/pq[*]
		SQLite:github.com/mattn/go-sqlite3[*]
		MS ADODB: github.com/mattn/go-adodb[*]
		ODBC: bitbucket.org/miquella/mgodbc[*]

	安装
		beedb支持go get方式安装，是完全按照Go Style的方式来实现的。
		go get github.com/astaxie/beedb

### 2. Cookie与Session
	import (
	"github.com/astaxie/session"
	_ "github.com/astaxie/session/providers/memory"
	)

	//全局变量session管理器
	var globalSessions *session.Manager


	//设置Cookie
	expiration := time.Now()
	expiration = expiration.AddDate(1, 0, 0)
	setCookie := http.Cookie{Name: "username", Value: "Zero", Expires: expiration}
	http.SetCookie(rw, &setCookie)

	//Go读取cookie
	getCookie, _ := req.Cookie("username")
	f.Fprint(rw, getCookie)
	//另一种循环读取所有Cookie的方式
	for _, getCookie := range req.Cookies() {
		f.Fprint(rw, getCookie.Name)
	}

	//启用Session
	sess := globalSessions.SessionStart(rw, req)
	//关闭Session
	//globalSessions.SessionDestroy(rw, req)
	//Session设置
	sess.Set("username", "Zero")

	//Session的值获取
	sess.Get("username")

### 3. Gin框架学习
	获取gin框架：
		go get github.com/gin-gonic/gin

```go
func main() {
	r := gin.Default()	//路有引擎
	r.GET("/Get", getMsg)  //定义get访问地址
	r.POST("/Post", postMsg)  //定义post访问地址
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

	r.Run("127.0.0.1:9090")  //如果不指定IP地址、端口号，默认为本地ip地址，8080端口，可直接简写端口号“:9090”
}

func test1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg": "返回信息",
		"data": "欢迎您！",
	})
}

//Post方式处理的请求
func postMsg(c *gin.Context) {
	name1 := c.DefaultPostForm("name", "Gin")
	fmt.Println(name1)
	name2, _ := c.GetPostForm("name")  //返回两个值
	fmt.Println(name2)
	c.String(http.StatusOK, "欢迎您的访问：%s", name1)
}

//Get方式处理的请求
func getMsg(c * gin.Context)  {
	name := c.Query("name")  //获取URL中的值
	//返回字符串类型
	//c.String(http.StatusOK, "欢迎您的访问：%s", name)
	//返回JSON数据类型
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg": "返回信息",
		"data": "欢迎您：" + name,
	})
}
```
