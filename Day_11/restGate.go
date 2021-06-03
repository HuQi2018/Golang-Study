/*
创建者：     Zero
创建时间：   2021/5/21
项目名称：   golang-study
*/
package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pjebs/restgate"
	"net/http"
)

/*
	RestGate：REST API端点安全认证中间件。
	优点：
		1、Go语言实现，可以与Gin框架无缝对接
		2、良好的设计，易于使用
	Gin框架使用RestGate
	1、添加依赖
		go get github.com/pjebs/restgate
	2、使用方法
		r.Use(authMiddleware()) //中间件
		restgate.New("X-Auth-Key", "X-Auth-Secret", restgate.Static,
					 restgate.Config{Key: []string{"12345", "gin"},Secret: []string{"secret", "gin_ok"}, })
*/

var mysqlDb *sql.DB

func init() {
	mysqlDb, _ = SqlDB()
}
func SqlDB() (*sql.DB, error) {
	//==============数据库连接方式1：==============
	DB_TYPE := "mysql"
	DB_HOST := "localhost"
	DB_PORT := "3306"
	DB_USER := "root"
	DB_NAME := "gomysql"
	DB_PASSWORD := "ok"
	openString := DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME
	db, err := sql.Open(DB_TYPE, openString)
	return db, err
	//==============数据库连接方式2：==============
	/*
		//parseTime:时间格式转换;
		// loc=Local解决数据库时间少8小时问题
		var err error
		db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/api_secure?charset=utf8&parseTime=true&loc=Local")
		if err != nil {
			log.Fatal("数据库打开出现了问题：", err)
			return db, err
		}
		// 尝试与数据库建立连接（校验dsn是否正确）
		err = db.Ping()
		if err != nil {
			log.Fatal("数据库连接出现了问题：", err)
			return db, err
		}
		return db, err
	*/
}

//验证方法
func authFunc1(c *gin.Context) {
	resData := HttpRes{
		Code:   http.StatusOK,
		Result: "验证通过！",
	}
	c.JSON(http.StatusOK, resData)
}

//验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1、验证  静态方法配置秘钥对
		//gate := restgate.New("X-Auth-Key", "X-Auth-Secret", restgate.Static, restgate.Config{
		//	Key: []string{"admin", "gin"},
		//	Secret: []string{"adminpw", "gin_ok"}, //秘钥组，上下对应，
		//	HTTPSProtectionOff: false,  //如果需要使用http则必须要设置为true
		//})
		//动态方法配置秘钥对
		//DB_TYPE := "mysql"
		//openString := "root:ok@tcp(localhost:3306)/gomysql"
		//mysqlDb, err := sql.Open(DB_TYPE, openString) //数据库连接
		//if err != nil {
		//	fmt.Println("数据库连接失败！", err)
		//}
		gate := restgate.New("X-Auth-Key", "X-Auth-Secret",
			restgate.Database, restgate.Config{
				DB:        mysqlDb,              //存储秘钥的数据库
				TableName: "users",              //存放秘钥对的表名称
				Key:       []string{"username"}, //key对应的列名
				Secret:    []string{"password"}, //secret对应的列名
			})
		nextCalled := false
		//下一个处理方法
		nextAdapter := func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			//继续往下执行
			c.Next()
		}
		//2、执行服务，如果都通过就执行nextAdapter方法
		//开启restgate
		gate.ServeHTTP(c.Writer, c.Request, nextAdapter)
		if !nextCalled {
			c.AbortWithStatus(401)
		}
	}
}
