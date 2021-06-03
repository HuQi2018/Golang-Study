package main

import (
	f "fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	//获取当前路径
	str, _ := os.Getwd()
	f.Println(str)

	//f.Println("Web：")
	//simpleWebFunc()
	//f.Println("")

	Bob := Human{"Bob", 39, "000-7777-XXX"}
	//任何实现了String方法的类型都能作为参数被fmt.Println调用
	f.Println("This Human is : ", Bob) //This Human is :  ❰Bob - 39 years -  ✆ 000-7777-XXX❱

	http.HandleFunc("/", sayhelloName) //设置访问的路由
	http.HandleFunc("/login", login)   //设置访问的路由
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
