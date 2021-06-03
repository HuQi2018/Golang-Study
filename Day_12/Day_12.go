package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//var sqlResponse HttpResponse

var sqlDb *sql.DB
var xormDb *xorm.Engine
var gormDb *gorm.DB

func init() {
	//普通数据库初始化
	//1、打开数据库
	//parseTime：时间格式转换（查询结果为时间时，是否自动解析为时间）；
	//Loc=Local:MySQL的时区设置 解决数据库时间少8小时问题
	DbType := "mysql"
	DbHost := "localhost"
	DbPort := "3306"
	DbUser := "root"
	DbName := "gomysql"
	DbPassword := "ok"
	setting := "?charset=utf8&parseTime=true&loc=Local"
	var err error
	openString := DbUser + ":" + DbPassword + "@tcp(" + DbHost + ":" + DbPort + ")/" + DbName + setting
	sqlDb, err = sql.Open(DbType, openString)
	if err != nil {
		fmt.Println("数据库打开出现了问题！", err)
		return
	}
	//2、测试与数据库连接的建立（校验连接是否正确）
	err = sqlDb.Ping()
	if err != nil {
		fmt.Println("数据库连接出现了问题！", err)
		return
	}
	//return db, err

	//ORM连接数据库初始化
	xormDb, err = xorm.NewEngine(DbType, openString) //创建数据库引擎
	if err != nil {
		fmt.Println("数据库连接失败！", err)
		return
	}
	//表同步，不存在表时则自动创建
	err = xormDb.Sync(new(Stu))
	if err != nil {
		fmt.Println("数据库同步失败！")
		return
	}

	//GORM连接数据库初始化
	setting = "?charset=utf8mb4&parseTime=true&loc=Local"
	openString = DbUser + ":" + DbPassword + "@tcp(" + DbHost + ":" + DbPort + ")/" + DbName + setting
	gormDb, err = gorm.Open(mysql.Open(openString), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败！")
		return
	}
	//表同步，不存在表时则自动创建
	err = gormDb.AutoMigrate(&Product{})
	if err != nil {
		fmt.Println("数据库同步失败!")
	}
}

func main() {
	r := gin.Default()
	r.POST("/sql/insertUser", insertData)   //添加数据
	r.GET("/sql/getUser", getUser)          //查询单条数据
	r.GET("/sql/getUsers", getMulData)      //查询多条数据
	r.PUT("/sql/updateUser", updateData)    //更新数据
	r.DELETE("/sql/deleteUser", deleteData) //删除数据

	r.POST("/orm/insertStu", xormInsertData)   //添加数据
	r.GET("/orm/getStu", xormGetData)          //查询单条数据
	r.GET("/orm/getStus", xormGetMulData)      //查询多条数据
	r.PUT("/orm/updateStu", xormUpdateData)    //更新数据
	r.DELETE("/orm/deleteStu", xormDeleteData) //删除数据

	r.POST("/gorm/insertPro", gormInsertData)   //添加数据
	r.GET("/gorm/getPro", gormGetData)          //查询单条数据
	r.GET("/gorm/getPros", gormGetMulData)      //查询多条数据
	r.PUT("/gorm/updatePro", gormUpdateData)    //更新数据
	r.DELETE("/gorm/deletePro", gormDeleteData) //删除数据

	err := r.Run(":9090")
	if err != nil {
		fmt.Println("服务运行失败！", err)
	}
}
