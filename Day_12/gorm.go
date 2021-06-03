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
gorm也是一款非常优秀的Go语言orm框架
1、获取框架包
go get -u gorm.io/gorm
go get gorm.io/driver/mysql
2、连接数据库
gorm.Open(mysql.Open(sqlStr), &gorm.Config{})
3、操作数据库
db.Create(p)  //p：结构体对象
db.Where(p.ID).Delete(P{})
db.Where("id=?", p.ID).Find(&p)
db.Model(&p).Where("number", p.Number).Updates(&p)
*/

//结构体名称为：Product，创建表的名称为：Products
type Product struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Number         string    `gorm:"unique" json:"number"`                       //商品唯一编号
	Category       string    `gorm:"type:varchar(256);not null" json:"category"` //商品类别
	Name           string    `gorm:"type:varchar(20);not null" json:"name"`      //商品名称
	MadeIn         string    `gorm:"type:varchar(128);not null" json:"made_in"`  //生产地
	ProductionTime time.Time `json:"production_time"`                            //生产时间
}

func gormDeleteData(c *gin.Context) {

	//方法捕获异常
	defer func() {
		err := recover()
		if err != nil {
			c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "错误！", err})
			return
		}
	}()

	number := c.Query("number")
	//1、先查询
	var count int64
	rs := gormDb.Model(Product{}).Where("number=?", number).Count(&count)
	if rs.Error != nil || count <= 0 {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "数据获取失败！", rs.Error})
		return
	}
	//2、再删除
	rs = gormDb.Where("number=?", number).Delete(&Product{})
	if rs.RowsAffected <= 0 || rs.Error != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "数据删除失败！", rs.Error})
		return
	}
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "数据删除成功！", rs.RowsAffected})
}

func gormUpdateData(c *gin.Context) {
	var product Product
	err := c.Bind(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "参数错误！", err})
		return
	}
	//1、先查询
	var count int64
	rs := gormDb.Model(Product{}).Where("number=?", product.Number).Count(&count)
	if rs.Error != nil || count <= 0 {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "数据获取失败！", rs.Error})
		return
	}
	//2、再更新
	rs = gormDb.Model(Product{}).Where("number=?", product.Number).Updates(&product)
	//无数据改变时更新行数也是0
	if rs.RowsAffected <= 0 || rs.Error != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "更新失败！", rs.Error})
		return
	}
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "更新成功！", rs.RowsAffected})
}

func gormGetMulData(c *gin.Context) {
	category := c.Query("category")
	products := make([]Product, 10)
	rs := gormDb.Where("category=?", category).Find(&products).Limit(10)
	if rs.Error != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "查询失败！", rs.Error})
		return
	}
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "查询成功", products})
}

func gormGetData(c *gin.Context) {
	number := c.Query("number")
	var product Product
	rs := gormDb.Where("number=?", number).First(&product)
	if rs.Error != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "查询失败！", rs.Error})
		return
	}
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "查询成功", product})
}

func gormInsertData(c *gin.Context) {
	var prod Product
	err := c.Bind(&prod)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "请求参数有误！", err})
		return
	}
	rs := gormDb.Create(&prod)
	if rs.RowsAffected <= 0 {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "数据添加失败！", rs.Error})
		return
	}
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "添加成功！", rs.RowsAffected})
}
