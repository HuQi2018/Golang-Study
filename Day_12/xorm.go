/*
创建者：     Zero
创建时间：   2021/5/24
项目名称：   golang-study
*/
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

/*
ORM，即Object-Relationl Mapping。
它的作用是在关系型数据库和对象之间作一个映射，这样我们再具体操作数据库的时候，就不需要再去和复杂的SQL语句打交道，
只要像平时操作对象一样操作它就可以了。
比较好的Go语言ORM包括：XORM、GORM。
XORM是一个简单而强大的Go语言ORM库，通过它可以使数据库操作非常简便。

XORM使用：
1、添加依赖：
go get github.com/go-xorm/xorm
2、连接数据库：
//创建引擎
engine, err := xorm.NewEngine(driverName, dataSourceName)
//连接  定义一个和表同步的结构体，并且自动同步结构体到数据库
type User struct{
	Id int64
	Name string
	Age int
}
err := engine.Sync2(new(User))
3、执行操作
engine.Get(&user)
engine.Insert(&user)
engine.ID(1).Uodate(&user)
engine.ID(1).Delete(&user)
*/

type Stu struct { //结构体名称与表的名称要一致
	Id      int64     `xorm:"pk autoincr" json:"id"` //主键自增
	StuNum  string    `xorm:"unique" json:"stu_num"`
	Name    string    `json:"name"`
	Age     int       `json:"age"`
	Created time.Time `xorm:"created" json:"created"`
	Updated time.Time `xorm:"updated" json:"updated"`
}

func xormDeleteData(c *gin.Context) {
	stuNum := c.Query("stu_num")
	//1、先查找
	var stus []Stu
	err := xormDb.Where("stu_num=?", stuNum).Find(&stus)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "请求参数有误！", err})
		return
	}
	//2、再删除
	affectedRow, err := xormDb.Where("stu_num=?", stuNum).Delete(&Stu{})
	if err != nil || affectedRow <= 0 {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "删除数据失败！", err})
		return
	}
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "删除成功！", affectedRow})
}

func xormUpdateData(c *gin.Context) {
	var stu Stu
	err := c.Bind(&stu)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "请求参数有误！", err})
		return
	}
	//1、先查找
	var stus []Stu
	err = xormDb.Where("stu_num=?", stu.StuNum).Find(&stus)
	if err != nil || len(stus) <= 0 {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "查询数据不存在！", err})
		return
	}
	//2、再更新
	affectedRow, err := xormDb.Where("stu_num=?", stu.StuNum).Update(&Stu{Name: stu.Name, Age: stu.Age})
	if err != nil || affectedRow <= 0 {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "更新失败！", err})
		return
	}
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "更新成功！", affectedRow})
}

func xormGetMulData(c *gin.Context) {
	name := c.Query("name")
	var stus []Stu
	err := xormDb.Where("name=?", name).And("age > 20").Limit(10, 0).Asc("age").Find(&stus)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "查询失败！", err})
		return
	}
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "查询成功！", stus})
}

func xormGetData(c *gin.Context) {
	stuNum := c.Query("stu_num")
	var stus []Stu
	err := xormDb.Where("stu_num=?", stuNum).Find(&stus)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "查询错误！", err})
		return
	}
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "查询成功！", stus})
}

func xormInsertData(c *gin.Context) {
	var stu Stu
	err := c.Bind(&stu)
	if err != nil || stu.StuNum == "" {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "参数错误！", err})
		return
	}
	affectedRow, err := xormDb.Insert(stu)
	if affectedRow <= 0 || err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "添加数据失败！", err})
		return
	}
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "添加数据成功！", string(affectedRow)})
}
