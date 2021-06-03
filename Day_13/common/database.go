/*
创建者：     Zero
创建时间：   2021/5/25
项目名称：   golang-study
*/
package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"golang-study/huqi/Day_13/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/url"
)

var MyDB *gorm.DB

func InitDB() *gorm.DB {
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	loc := viper.GetString("datasource.loc")

	//setting := "?charset=" + charset + "&parseTime=true&loc=" + loc
	//openString := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + setting
	//fmt.Println(openString)
	//fmt.Printf("%s::%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
	//	username, password, host, port, database, charset, url.QueryEscape(loc))
	sqlStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username, password, host, port, database, charset, url.QueryEscape(loc))
	db, err := gorm.Open(mysql.Open(sqlStr), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败：" + err.Error())
	}
	MyDB = db
	//表同步，不存在表时则自动创建
	err = MyDB.AutoMigrate(&model.UserBase{})
	if err != nil {
		fmt.Println("数据库同步失败!")
	}
	return MyDB
}
