/*
创建者：     Zero
创建时间：   2021/5/20
项目名称：   golang-study
*/
package main

/*
Swagger是一个规范和完整的框架，用于生成、描述、调用和可视化
RESTful风格的Web服务
	作用：
		（1）接口的文档在线自动生成
		（2）功能测试
	使用：
		（1）添加依赖
			go get github.com/swaggo/swag/cmd/swag
			go get github.com/swaggo/gin-swagger
			go get github.com/swaggo/gin-swagger/swaggerFiles
		（2）给Handler对应方法添加注释
			@Tags：说明该方法的作用
			@Summary：登录
			@Description：这个API详细的描述
			@Accept：表示该请求的请求类型
			@Produce：返回数据类型
			@Param：参数，表示需要传递到服务器端的参数    变量名 get、post-》query 值的类型 是否不允许为空 备注
			@Success：成功返回给客户端的信息
			@Router：路由信息
*/

//因swagger必须在main.go文件内，所以具体内容请查看main.go文件。
