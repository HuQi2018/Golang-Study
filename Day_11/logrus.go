/*
创建者：     Zero
创建时间：   2021/5/21
项目名称：   golang-study
*/
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	lfshook "github.com/rifflock/Lfshook"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

/*
使用Logrus来记录日志
Logrus是一个结构化、可插拔的Go日志框架，完全兼容官方log库接口。
功能强大的同时，Logrus具有高度的灵活性，他提供了自定义插件的功能，有TEXT与JSON两种可选的日志输出格式。
Logrus还支持Field机制和可扩展的HOOK机制。
它鼓励用户通过Field机制进行精细化的、结构化的日志记录，允许用户通过hook的方式将日志分发到任意地方。
许多著名开源项目，如docker、prometheus等都是使用Logrus来记录日志。

使用Logrus
	1、添加依赖
	go get github.com/sirupsen/logrus
	2、示例代码
	log.WithFields(log.Fields{"name": "张三",}).Info("这里是日志信息")

日志切割：
1、通过"github.com/lestrrat-go/file-rotatelogs"完成文件切割。
	logWriter, err := rotatelogs.New(
		//分割后的文件名称
		fileName + ".%Y%m%d.log",
		//生成软链，指向最新日志文件
		//如果文件的软链接存在，则会提示错误（并非panic错误）
		rotatelogs.WithLinkName(fileName),
		//设置最大保存时间（7天）
		rotatelogs.WithMaxAge(7 * 24 * time.Hour),
		//设置日志切割时间间隔（1天）
		rotatelogs.WithRotationTime(24 * time.Hour),
	)
2、通过"github.com/rifflock/lfshook"完成log文件的hook机制。
writeMap := lfshook.WriterMap{
	logrus.InfoLevel: logWriter,
	logrus.FatalLevel: logWriter
}
*/

var log = logrus.New() //创建一个Log示例
func initLogrus() error {
	log.Formatter = &logrus.JSONFormatter{}                                                          //设置为json格式的日志
	file, err := os.OpenFile("./huqi/Day_11/gin_log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) //创建一个Log日志文件
	if err != nil {
		fmt.Println("创建文件/打开文件失败！")
		return err
	}
	log.Out = file               //设置log的默认文件输出
	gin.SetMode(gin.ReleaseMode) //发布版本
	gin.DefaultWriter = log.Out  //将gin框架自己记录的日志也输出到日志信息中
	log.Level = logrus.InfoLevel //设置日志级别
	return nil
}

func logrusFunc1(c *gin.Context) {
	log.WithFields(logrus.Fields{ //设置日志内容
		"url":    c.Request.RequestURI,
		"method": c.Request.Method,
		"params": c.Query("name"),
		"IP":     c.ClientIP(),
	}).Info()
	resData := HttpRes{
		Code:   http.StatusOK,
		Result: "响应成功！",
	}
	c.JSON(http.StatusOK, resData)
}

//log切割日志的中间件
func LogMiddleware() gin.HandlerFunc {
	//日志文件
	//fileName := path.Join(logFilePath, logFileName)
	fileName := "huqi/Day_11/system_log.log"
	//写入文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
	}
	//实例化
	logger := logrus.New()
	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	//设置输出
	logger.Out = file
	//设置rotatelogs，实现文件分割
	logWriter, err := rotatelogs.New(
		//分割后的文件名称
		fileName+"%Y%m%d.log",
		//生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		//设置最大保存时间（7天）
		rotatelogs.WithMaxAge(7*24*time.Hour), //以hour为单位的整数
		//设置日志切割时间间隔（1天）
		rotatelogs.WithRotationTime(1*time.Hour),
	)
	if err != nil {
		fmt.Println("日志切割失败！", err)
	}
	//hook机制的设置
	writerMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	logger.AddHook(lfshook.NewHook(writerMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))
	return func(c *gin.Context) {
		c.Next()
		method := c.Request.Method
		reqUrl := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		logger.WithFields(logrus.Fields{
			"status_code": statusCode,
			"client_ip":   clientIP,
			"req_method":  method,
			"req_uri":     reqUrl,
		}).Info()
	}
}

func logrusFunc2(c *gin.Context) {
	resData := HttpRes{
		Code:   http.StatusOK,
		Result: "响应成功！",
	}
	c.JSON(http.StatusOK, resData)
}
