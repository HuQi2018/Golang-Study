package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//这段代码不再使用 http.HandleFunc 函数，取而代之的是直接调用 http.Handle 并传入我们自定义的 http.Handler 实现。
//初学者有一点需要特别注意，即 Go 语言是一门大小写敏感的语言（否则无法通过首字母大小写区分一个对象是公开的还是私有的）。
//因此，想要实现http.Handler 接口，方法名称必须连大小写也保持一致，即这里的方法名称必须是 ServeHTTP 而不可以是 ServeHttp。

func selfHandle() {
	http.Handle("/", &helloHandler{})

	log.Println("Starting HTTP server...")
	log.Fatal(http.ListenAndServe(":4000", nil))
}

/*
启动服务器后，如果我们访问 http://localhost:4000 会收到和之前一样的结果，但如果尝试访问 http://localhost:4000/timeout 则不会收到任何消息：
这是因为我们的执行函数在休眠 2 秒后被 http.Server 对象认为已经超时，提前关闭了与客户端之间的连接，因此无论执行函数后面向响应体写入任何东西都不会有任何作用。
*/
func selfTimeOut() {
	mux := http.NewServeMux()
	mux.Handle("/", &helloHandler{})
	mux.HandleFunc("/timeout", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Write([]byte("Timeout"))
	})

	server := &http.Server{
		Addr:         ":4000",
		Handler:      mux,
		WriteTimeout: 2 * time.Second,
	}
	log.Println("Starting HTTP server...")
	log.Fatal(server.ListenAndServe())
}

/*
通过结合捕捉系统信号（Signal）、goroutine 和管道（Channel）来实现服务器的优雅停止
这段代码通过捕捉 os.Interrupt 信号（Ctrl+C）然后调用 server.Shutdown 方法告知服务器应停止接受新的请求并在处理完当前已接受的请求后关闭服务器。
为了与普通错误相区别，标准库提供了一个特定的错误类型 http.ErrServerClosed，我们可以在代码中通过判断是否为该错误类型来确定服务器是正常关闭的还是意外关闭的。
*/
func closeRequest() {
	mux := http.NewServeMux()
	mux.Handle("/", &helloHandler{})

	server := &http.Server{
		Addr:    ":4000",
		Handler: mux,
	}

	// 创建系统信号接收器
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit

		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal("Shutdown server:", err)
		}
	}()

	log.Println("Starting HTTP server...")
	err := server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			log.Print("Server closed under request")
		} else {
			log.Fatal("Server closed unexpected")
		}
	}
}

type helloHandler struct{}

//自定义实现http.Handler
func (*helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

/*
template.New 的作用就是根据用户给定的名称创建一个模板对象，本例中我们使用了 “test” 字符串作为这个模板对象的名称。
另外，由于 template.New 函数会直接返回一个 *template.Template 对象，因此可以直接链式操作调用该对象的 Parse 方法
template.Parse 方法接受一个 string 类型的参数，即文本模板的内容，然后对内容进行解析并返回解析过程中发生的任何错误。
本例中，我们使用了没有任何模板语法的 “Hello world!” 字符串，同时获得了两个返回值。第一个返回值依旧是一个 *template.Template 对象，
此时该对象已经包含了模板解析后的数据结构。第二个返回值便是在解析过程中可能出现的错误，这要求我们对该错误进行检查判断。
如果模板解析过程没有产生任何错误则表示模板可以被用于渲染了，template.Execute 就是用于渲染模板的方法，该方法接受两个参数：
输出对象和指定数据对象（或根对象）。简单起见，本例中我们只使用到了第一个参数，即输出对象。
凡是实现了 io.Writer 接口的实例均可以作为输出对象，这在 Go 语言中是非常常见的一种编码模式。

template.Execute 方法的第二个参数类型为 interface{}，也就是说可以传入任何类型。
{{.}} 点操作符默认指向的是根对象
而如果根对象为一个复合类型，那么点操作符所代表的也就是这个复合类型。
http://localhost:4000/?sku=1122334&name=phone&unitPrice=649.99&quantity=833
Inventory
SKU: 1122334
Name: phone
UnitPrice: 649.99
Quantity: 833
*/
//在模板中渲染复杂对象
type Inventory struct {
	SKU       string
	Name      string
	UnitPrice float64
	Quantity  int64
}

// 在模板中调用结构的方法
// Subtotal 根据单价和数量计算出总价值
func (i *Inventory) Subtotal() float64 {
	return i.UnitPrice * float64(i.Quantity)
}
func templateHandle() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		/*
			html/template 与 text/template 的关联与不同
			html/template 本身是一个 text/template 包的一层封装，并在此基础上专注于提供安全保障。
			作为使用者来说，最直观的变化就是对所有的文本变量都进行了转义处理。
			反转义  html/template
			{{.content | safe}}

			修改分隔符
			tmpl, err := template.New("test").Delims("[[", "]]").Parse(`[[.content]]`)
			我们通过 Delims 方法将它们分别修改为方括号 [[ 和 ]]。
		*/
		// 创建模板对象并添加自定义模板函数 .Funcs
		// 创建模板对象并解析模板内容
		tmpl, err := template.New("test").Funcs(template.FuncMap{
			"add": func(a, b int) int {
				return a + b
			},
			"add2": func(a int) int {
				return a + 2
			},
			"join": strings.Join,
		}).Parse(`Val is {{.val}} !
{{/* 模板中使用自定义函数 */}}
result: {{add 1 2}}
{{/* 模板中使用管道 */}}
result: {{add2 0 | add2 | add2}}
result: {{add 1 3 | add 2 | add 2}}
{{/* 模板复用  使用到了模板复用最核心的概念：定义、使用和传参。
通过 Funcs 方法添加了名为 join 模板函数，其实际上就是调用 strings.Join
通过 define "<名称>" 的语法定义了一个非常简单的局部模板，即以根对象 . 作为参数调用 join 模板函数
通过 template "<名称>" <参数> 的语法，调用名为 list 的局部模板，并将 .names 作为参数传递进去（传递的参数会成为局部模板的根对象）
Names: Alice, Bob, Cindy, David
*/}}
{{define "list"}}
    {{join . ", "}}
{{end}}
Names: {{template "list" .names}}`)
		// 使用模板文件
		//tmpl, err := template.ParseFiles("huqi/Day_07/template_local.tmpl")
		if err != nil {
			//log.Fatalf("Parse: %v", err)
			fmt.Fprintf(w, "Parse: %v", err)
			return
		}

		//在模板中渲染变量
		// 获取 URL 参数的值
		val := r.URL.Query().Get("val")
		// 根据 URL 查询参数的值创建 Inventory 实例
		inventory := &Inventory{
			SKU:  r.URL.Query().Get("sku"),
			Name: r.URL.Query().Get("name"),
		}
		// 注意：为了简化代码逻辑，这里并没有进行错误处理
		inventory.UnitPrice, _ = strconv.ParseFloat(r.URL.Query().Get("unitPrice"), 64)
		inventory.Quantity, _ = strconv.ParseInt(r.URL.Query().Get("quantity"), 10, 64)

		//在模板中使用条件判断（if 语句）
		x, _ := strconv.ParseInt(r.URL.Query().Get("x"), 10, 64)
		y, _ := strconv.ParseInt(r.URL.Query().Get("y"), 10, 64)
		// 当 y 不为 0 时进行除法运算
		yIsZero := y == 0
		result := 0.0
		if !yIsZero {
			result = float64(x) / float64(y)
		}

		//使用 map 类型作为模板根对象
		// 调用模板对象的渲染方法，并创建一个 map[string]interface{} 类型的临时变量作为根对象
		err = tmpl.Execute(w, map[string]interface{}{
			"val":       val,
			"Inventory": inventory,
			"yIsZero":   yIsZero,
			"result":    result,
			"arrLists": []string{
				"Alice",
				"Bob",
				"Carol",
				"David",
			},
			"Numbers": []int{1, 3, 5, 7},
			"names":   []string{"Alice", "Bob", "Cindy", "David"},
		})
		if err != nil {
			fmt.Fprintf(w, "Execute: %v", err)
			return
		}
	})
	//获取 “/?val=123” 中的 “val” 的值，并返回给客户端
	http.HandleFunc("/get_var", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.Query().Get("val")))
	})
	log.Println("Starting HTTP server...")
	log.Fatal(http.ListenAndServe("localhost:4000", nil))
}
