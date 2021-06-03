package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）  解析参数，默认是不会解析的
	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

type User struct {
	Id       int
	UserName string
	RealName string
	Age      int
	Password string
	UserCard string
	Email    string
	Mobile   string
}
type Register1 struct {
	User
	Fruit    string
	Interest string
	Message  string
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("huqi/Day_06/login.gtpl")
		t.Execute(w, token)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		token := r.Form.Get("token")
		if token != "" {
			//验证token的合法性
		} else {
			//不存在token报错
		}
		fmt.Println("username length:", len(r.Form["username"][0]))
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) //输出到服务器端
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		template.HTMLEscape(w, []byte(r.Form.Get("username"))) //输出到客户端
	}
}

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	t, _ := template.ParseFiles("huqi/Day_06/upload.gtpl")
	message := ""
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		message = fmt.Sprintf("%x", h.Sum(nil))

	} else {
		//处理的最大上传数据大小
		//上传的文件存储在maxMemory大小的内存里面，如果文件大小超过了maxMemory，那么剩下的部分将存储在系统的临时文件中
		//上传文件主要三步处理：
		//表单中增加enctype="multipart/form-data"
		//服务端调用r.ParseMultipartForm,把上传的文件存储在内存和临时文件中
		//使用r.FormFile获取文件句柄，然后对文件进行存储等处理。
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		//fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("huqi/Day_06/upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		message = "上传成功！"
	}
	t.Execute(w, message)
}
func login2(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)                      //获取请求的方法
	t, _ := template.ParseFiles("huqi/Day_06/login.gtpl") //解析模板
	//t.Execute(w, nil)                            //渲染模板，并发送
	//if r.Method == "POST" {
	//请求的是登陆数据，那么执行登陆的逻辑判断
	//解析表单
	//Request本身也提供了FormValue()函数来获取用户提交的参数。如r.Form["username"]也可写成r.FormValue("username")。
	//调用r.FormValue时会自动调用r.ParseForm，所以不必提前调用。
	//r.FormValue只会返回同名参数中的第一个，若参数不存在则返回空字符串。
	r.ParseForm()
	user := User{}
	user.UserName = r.Form.Get("username")
	user.Password = r.Form.Get("password")
	user.Age, _ = strconv.Atoi(r.Form.Get("age"))
	user.Email = r.Form.Get("email")
	user.RealName = r.Form.Get("realname")
	user.Mobile = r.Form.Get("mobile")
	user.UserCard = r.Form.Get("usercard")
	register := Register1{User: user}
	register.Fruit = r.Form.Get("fruit")
	register.Interest = r.Form.Get("interest")

	_, mesg := formValidFunc(r)
	register.Message = mesg
	//if ok{
	//	register.Message = mesg
	//}

	fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) //输出到服务器端
	fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
	//template.HTMLEscape(w, []byte(r.Form.Get("username"))) //输出到客户端 //XSS预防
	t.Execute(w, mesg)
	//}
}

func formValidFunc(r *http.Request) (tag bool, message string) {
	tag = true
	r.ParseForm()
	message = ""
	//为空的处理
	if len(r.Form.Get("username")) == 0 {
		message = "用户名不能为空！"
	}
	//判断正整数
	getint, err := strconv.Atoi(r.Form.Get("age"))
	if err != nil {
		//数字转化出错了，那么可能就不是数字
		return false, "年龄请输入数字！"
	}
	//接下来就可以判断这个数字的大小范围了
	if getint > 100 {
		return false, "年龄超过100岁，不符合要求！"
	}
	//正则匹配的方式
	if m, _ := regexp.MatchString("^[0-9]+$", r.Form.Get("age")); !m {
		return false, "年龄数据不符合要求！"
	}
	//中文正则
	if m, _ := regexp.MatchString("^\\p{Han}+$", r.Form.Get("realname")); !m {
		return false, "请输入中文名字！"
	}
	//英文正则
	if m, _ := regexp.MatchString("^[a-zA-Z]+$", r.Form.Get("username")); !m {
		return false, "用户名只能包含英文！"
	}
	//邮箱地址判断
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, r.Form.Get("email")); !m {
		return false, "请输入正确的邮箱地址！"
	}
	//手机号码判断
	if m, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, r.Form.Get("mobile")); !m {
		return false, "请输入正确的手机号码！"
	}
	//单选按钮、下拉菜单验证
	slice := []string{"apple", "pear", "banane"}
	if !in(r.Form.Get("fruit"), slice) {
		return false, "水果数据非法！"
	}
	//日期和时间验证
	//t := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	//fmt.Printf("Go launched at %s\n", t.Local())
	//身份证号码验证
	//验证15位身份证，15位的是全部数字
	if m, _ := regexp.MatchString(`^(\d{15})$`, r.Form.Get("usercard")); !m {
		return false, "身份证号错误！"
	}
	//验证18位身份证，18位前17位为数字，最后一位是校验位，可能为数字或字符X。
	if m, _ := regexp.MatchString(`^(\d{17})([0-9]|X)$`, r.Form.Get("usercard")); !m {
		return false, "身份证号错误！"
	}
	return true, ""
}

func in(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}
