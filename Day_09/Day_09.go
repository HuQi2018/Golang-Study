package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"strconv"
	"time"
)

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
	server02 := &http.Server{
		Addr:         ":9092",
		Handler:      router02(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	//开启服务
	g.Go(func() error { //开启服务程序1
		return server01.ListenAndServe()
	})
	g.Go(func() error {
		return server02.ListenAndServe()
	})
	if err := g.Wait(); err != nil {
		fmt.Println("执行失败：", err)
	}

	//r := gin.Default()       //获取路有引擎
	//r.GET("/Get", getMsg)    //定义get访问地址
	//r.Run(":9090") //如果不指定IP地址、端口号，默认为本地ip地址，8080端口，可直接简写端口号“:9090”
}

func router02() http.Handler {

	r := gin.Default() //获取路有引擎

	//获取外部地址的内容然后显示
	r.GET("/getBaiDuContent", func(c *gin.Context) {
		url := "https://mp.weixin.qq.com/s/HgoLvNF5p4kZVQn0YDG5jA"
		res, err := http.Get(url)
		if err != nil || res.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable) //应答Client
			return
		}
		body := res.Body
		contentLength := res.ContentLength
		contentType := res.Header.Get("Content-Type")
		//数据写入响应体
		c.DataFromReader(http.StatusOK, contentLength, contentType, body, nil)
	})

	//Gin框架的多形式渲染
	r.GET("/requestContent", func(c *gin.Context) {
		//返回JSON  前端会将<、>编码
		//c.JSON(http.StatusOK, gin.H{"html": "<b>Hello Gin框架</bin>"})
		//原样输出html
		//c.PureJSON(http.StatusOK, gin.H{"html": "<b>Hello Gin框架</bin>"})
		//返回yaml形式（yaml渲染）
		//c.YAML(http.StatusOK, gin.H{"message": "Gin框架的多形式渲染", "status": http.StatusOK})
		//输出xml形式（XML渲染）
		type Message struct {
			Name string
			Msg  string
			Age  int
		}
		data := Message{Name: "Gin框架", Msg: "Hello", Age: 123}
		c.XML(http.StatusOK, data)
	})

	//Gin框架响应文件
	//不存在是时返回404错误
	//必须进行限制，否则将出现安全问题，返回代码内容
	r.GET("/getFile", func(c *gin.Context) {
		filePath := "huqi/Day_09/Download/"
		fileName := filePath + c.Query("file")
		c.File(fileName)
	})

	//单个文件上传
	r.POST("/singleFileUpload", func(c *gin.Context) {
		//文件对应的key（Post方法）
		file, err := c.FormFile("file")
		if file == nil {
			c.String(http.StatusBadRequest, "不存在上传文件！")
			return
		}
		if err != nil {
			c.String(http.StatusBadRequest, "文件上传错误！")
		}
		fileUploadPath := "huqi/Day_09/Upload/"
		//存储文件
		err = c.SaveUploadedFile(file, fileUploadPath+file.Filename)
		if err != nil {
			c.String(http.StatusBadRequest, "文件上传错误！")
		}
		c.String(http.StatusOK, fmt.Sprintf("'%s' 上传完成！", file.Filename))
	})

	//多文件上传
	r.POST("/multiFileUpload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if form == nil {
			c.String(http.StatusBadRequest, "不存在提交表单！")
			return
		}
		if err != nil {
			c.String(http.StatusBadRequest, "上传文件错误！")
		}
		files := form.File["file_key"] //每个文件对应的key： file_key
		fileUploadPath := "huqi/Day_09/Upload/"
		for _, file := range files {
			fmt.Println("文件：" + file.Filename + "\t上传完成！")
			//上传文件至指定目录
			err := c.SaveUploadedFile(file, fileUploadPath+file.Filename)
			if err != nil {
				c.String(http.StatusBadRequest, "上传文件错误！"+file.Filename)
			}
		}
		c.String(http.StatusOK, fmt.Sprintf("%d 文件上传完成！", len(files)))
	})

	//Gin框架提供了快速登录验证中间件，可以完成登录的验证
	//使用gin.BasicAuth()中间件
	r.Use(AuthMiddleware())
	r.GET("/login", func(c *gin.Context) {
		//获取中间件的数据
		user := c.MustGet(gin.AuthUserKey).(string)
		c.JSON(http.StatusOK, user+"   登陆成功！")
	})

	v1 := r.Group("v1")
	{ //路由分组1（1级路径） 注意另起一行  也可不加大括号
		r := v1.Group("/user") //路由分组2（2级路径）
		r.GET("login", func(c *gin.Context) {
			//获取中间件的数据
			user := c.MustGet(gin.AuthUserKey).(string)
			c.JSON(http.StatusOK, user+"   登陆成功！"+c.Request.URL.Path)
		}) //响应请求：/v1/user/login
		r2 := r.Group("/showInfo")   //路由分组3（3级路径）
		r2.GET("/abstract", postMsg) //响应请求：/v1/user/showInfo/abstract
		r2.GET("/detail", test1)     //响应请求：/v1/user/showInfo/detail
	}

	//同步方法 调用一旦开始，调用者必须等到方法调用返回后，才能继续后续的行为
	//异步方法 调用更像一个消息传递，一旦开始，方法调用就会立即返回，调用者就可以继续后续的操作。而异步方法通常会在另一个go程（goroutine）过程，
	//不会阻碍调用者的工作。
	//可以在中间件或处理程序中启动新的Go程（goroutines）
	//特别注意：需要使用上下文的副本
	//同步请求，一条路执行完，才返回
	r.GET("/sync", func(c *gin.Context) {
		sync(c)
		c.JSON(http.StatusOK, ">>>主程序（主go程）同步已经执行<<<")
	})
	//异步请求，使用go开辟线程，不用等待线程执行完，直接返回
	r.GET("/async", func(c *gin.Context) {

		for i := 0; i < 6; i++ {
			cCp := c.Copy()
			go async(cCp, i)
		}
		c.JSON(http.StatusOK, ">>>主程序（主go程）同步已经执行<<<"+c.Request.URL.Path)

	})

	//是设置前端页面目录
	r.LoadHTMLGlob("huqi/Day_09/templates/*")
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	//自定义中间件
	//Gin中间件，对路由到来的数据先进行预处理，包括数据加载、请求验证（过滤）等。
	//gin.Default() //默认路由引擎，包括Logger and Recovery middleware
	//r := gin.New() //没有任何中间件的路由引擎
	r.Use(Middleware()) //使用添加中间件
	r.GET("/middleware", func(c *gin.Context) {
		fmt.Println("服务端执行开始。。。。")
		name := c.Query("name")
		ageStr := c.Query("age")
		age, _ := strconv.Atoi(ageStr)
		log.Println(name, age)
		rs := struct {
			Name string `json:"name"` //加入标签，将键值首字母小写
			Age  int    `json:"age"`
		}{name, age}
		c.JSON(http.StatusOK, rs)
	})
	return r
}

func router01() http.Handler {

	r := gin.Default()       //获取路有引擎
	r.GET("/Get", getMsg)    //定义get访问地址
	r.POST("/Post", postMsg) //定义post访问地址
	r.GET("/test1", test1)
	//外部重定向 可以通过Redirect跳转到外部页面
	//http.StatusMovedPermanently为状态码301 永久移动    请求的页面已永久跳转到新的url
	//第二个参数为跳转的外部地址
	r.GET("/redictToBaiDu", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com")
	})
	//内部路由重定向 通过c.Request.URL.Path 设置跳转的指定的路径
	//通过HandleContext函数   参数会连带着传给对应的地址
	r.GET("/redictToLocal", func(c *gin.Context) {
		// 指定重定向的URL 通过HandleContext进行重定向到test2 页面显示json数据
		c.Request.URL.Path = "/Get"
		r.HandleContext(c)
	})
	return r

}

//异步执行
func async(cCp *gin.Context, i int) {
	fmt.Println("第" + strconv.Itoa(i) + "个go程开始执行：" + cCp.Request.URL.Path)
	time.Sleep(time.Duration(i) * time.Second)
	fmt.Println("第" + strconv.Itoa(i) + "个go程执行结束++++")
}

//同步执行
func sync(c *gin.Context) {
	fmt.Println("开始执行同步任务：" + c.Request.URL.Path)
	time.Sleep(time.Second * 3)
	fmt.Println("同步任务执行完成！")
}

//定义登录中间件
func AuthMiddleware() gin.HandlerFunc {
	//初始化用户，静态添加
	//gin.Accounts是map[string]string类型
	accounts := gin.Accounts{
		"admin": "password",
	}
	//动态添加用户
	accounts["Golang"] = "123456"
	accounts["Gin"] = "789abc"
	auth := gin.BasicAuth(accounts)
	return auth
}

//自定义中间件
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("中间件开始执行=====")
		name := c.Query("name")
		ageStr := c.Query("age")
		age, err := strconv.Atoi(ageStr) //string-》int
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "年龄不是整数！")
			return
		}
		if age < 0 || age > 100 {
			c.AbortWithStatusJSON(http.StatusBadRequest, "年龄数据非法！")
			return
		}
		if len(name) < 6 || len(name) > 12 {
			c.AbortWithStatusJSON(http.StatusBadRequest, "用户名只能是6-12位！")
			return
		}
		c.Next() //执行后续操作
		//fmt.Println(name, age)
	}
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
