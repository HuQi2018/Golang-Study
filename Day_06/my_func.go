package main

import (
	f "fmt"
	"html/template"
	"net/http"
	"strconv"
)

func simpleWebFunc() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/login", index)
	http.ListenAndServe("localhost:8080", nil)
}
func index(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	//f.Fprint(rw, "Hello World!")
	message := ""
	if req.Method == "POST" {
		req.ParseForm()
		if cheackUserName(req.Form["username"][0]) {
			f.Fprint(rw, "用户名：", req.Form["username"][0], "密码：", req.Form["password"][0])
			return
		} else {
			message = "用户名长度不符合要求！"
		}
	}
	t, err := template.ParseFiles("huqi/Day_06/login.gtpl")
	if err != nil {
		f.Println(err)
	}
	t.Execute(rw, message)
}
func cheackUserName(username string) bool {
	if len(username) >= 6 && len(username) <= 16 {
		return true
	}
	return false
}
func handler(rw http.ResponseWriter, req *http.Request) {
	//解析用户请求的数据
	req.ParseForm()
	//判断用户是否有请求参数
	if len(req.Form["name"]) > 0 {
		f.Fprint(rw, "Hello ", req.Form["name"][0])
	} else {
		f.Fprint(rw, "Hello World!")
	}
}

type Human struct {
	name  string
	age   int
	phone string
}

// 通过这个方法 Human 实现了 fmt.Stringer
func (h Human) String() string {
	return "❰" + h.name + " - " + strconv.Itoa(h.age) + " years -  ✆ " + h.phone + "❱"
}
