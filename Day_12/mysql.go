/*
创建者：     Zero
创建时间：   2021/5/24
项目名称：   golang-study
*/
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserInfo struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	//Age        int    `json:"age"`
	//CreateTime string `json:"create_time"`
}

type HttpResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func insertData(c *gin.Context) {
	var user1 UserInfo
	//绑定数据
	err := c.Bind(&user1)
	if err != nil {
		fmt.Println("绑定失败！", err)
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "参数错误，绑定失败！", "Error!"})
		return
	}
	sqlStr := "insert into users(id, username, password) values (?,?,?)"
	rs, err := sqlDb.Exec(sqlStr, user1.Id, user1.Username, user1.Password)
	if err != nil {
		fmt.Printf("数据添加失败！%v\n", err)
		c.JSON(http.StatusOK, HttpResponse{http.StatusBadRequest, "数据添加失败！", err.Error()})
		return
	}
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "数据添加成功！", "Success！"})
	fmt.Println(rs.LastInsertId()) //打印添加后的userid
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getMulData(c *gin.Context) {
	name := c.Query("username")
	sqlStr := "select id,password from users where username=?"
	rows, err := sqlDb.Query(sqlStr, name)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "查询失败！", err.Error()})
		return
	}
	defer rows.Close()
	resUser := make([]UserInfo, 0)
	for rows.Next() {
		var user UserInfo
		rows.Scan(&user.Id, &user.Password)
		user.Username = name
		resUser = append(resUser, user)
	}
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "查询成功！", resUser})
}

func getUser(c *gin.Context) {
	name := c.Query("username")
	sqlStr := "select id,password from users where username=?"
	var user1 UserInfo
	//此处scan的接收个数和顺序必须与sql中的语句查出的数据一致   QueryRow只返回一行数据
	err := sqlDb.QueryRow(sqlStr, name).Scan(&user1.Id, &user1.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "查询失败！", err.Error()})
		return
	}
	user1.Username = name
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "查询成功！", user1})
}

func updateData(c *gin.Context) {
	var user UserInfo
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "参数错误！", err})
		return
	}
	//1、先查询
	var count int
	err = sqlDb.QueryRow("select count(*) from users where id=?", user.Id).Scan(&count)
	if count <= 0 || err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "要更新的数据不存在！", "error"})
		return
	}
	//再更新
	sqlStr := "update users set username=?, password=? where id=?"
	res, err := sqlDb.Exec(sqlStr, user.Username, user.Password, user.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "更新失败！", err})
		return
	}
	id, err := res.LastInsertId()
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "更新成功！", id})
	fmt.Println(id, err)
}

func deleteData(c *gin.Context) {
	id := c.Query("id")
	var count int
	//1、先查询
	sqlStr := "select count(*) from users where id=?"
	err := sqlDb.QueryRow(sqlStr, id).Scan(&count)
	if count <= 0 || err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "删除的数据不存在！", "error"})
		return
	}
	//2、再删除
	delStr := "delete from users where id=?"
	rs, err := sqlDb.Exec(delStr, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "删除的数据失败！", err})
		return
	}
	deleteId, err := rs.LastInsertId()
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "删除数据成功！", deleteId})
}
