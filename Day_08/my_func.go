package main

import (
	f "fmt"
	//"github.com/astaxie/beego/session"
	"github.com/astaxie/session"
	_ "github.com/astaxie/session/providers/memory"
	"net/http"
	"time"
)

type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string
	MaxAge     int
	Secure     bool
	HttpOnly   bool
	Raw        string
	Unparsed   []string // Raw text of unparsed attribute-value pairs
}

//全局session管理器
var globalSessions *session.Manager

//然后在init函数中初始化
func init() {
	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC()
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:8080", nil)
}
func handler(rw http.ResponseWriter, req *http.Request) {
	//设置Cookie
	expiration := time.Now()
	expiration = expiration.AddDate(1, 0, 0)
	setCookie := http.Cookie{Name: "username", Value: "Zero", Expires: expiration}
	http.SetCookie(rw, &setCookie)

	//Go读取cookie
	getCookie, _ := req.Cookie("username")
	f.Fprint(rw, getCookie)
	//另一种循环读取所有Cookie的方式
	for _, getCookie := range req.Cookies() {
		f.Fprint(rw, getCookie.Name)
	}

	//启用Session
	sess := globalSessions.SessionStart(rw, req)
	//关闭Session
	//globalSessions.SessionDestroy(rw, req)
	//Session设置
	sess.Set("username", "Zero")

	//Session的值获取
	sess.Get("username")

	//解析用户请求的数据
	req.ParseForm()
	//判断用户是否有请求参数
	if len(req.Form["name"]) > 0 {
		f.Fprint(rw, "Hello ", req.Form["name"][0])
	} else {
		f.Fprint(rw, "Hello World!")
	}

	//重定向
	//http.Redirect(w, r, "/", 302)
}
