package main

import (
	"errors"
	f "fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

func Go() {
	f.Println("Go并发函数！")
}

//普通通道
func concurrencyFunc() {
	//make创建的channel默认为双向通道，可存可取

	channel1 := make(chan bool)
	go Go()
	go func() {
		f.Println("Go匿名并发函数！")
		channel1 <- true
		//close(channel1)
	}()
	//防止主程序提前关闭，导致未有输出
	//time.Sleep(2 * time.Second)
	//使用管道让程序等待线程结束消息，再关闭主程序
	<-channel1 //取
	/*具体查看管道信息，在对channel类型进行迭代操作时，要注意在摸个地方关闭channel，
	  否则会产生死锁错误，如上匿名函数中在程序执行的最后调用close方法关闭channel*/
	//for v := range channel1{
	//	f.Println(v)
	//}
	/*
		有缓存则为同步阻塞，
		无缓存则为异步的，取的操作先于放
	*/
	channel2 := make(chan bool, 1) //设置缓存为1，不管你是否读出，都会结束
	go func() {
		f.Println("有缓存读！")
		<-channel2
	}()
	channel2 <- true
}

//解决并发异常1
func concurrencyFunc1() {
	//设置进程并发数
	runtime.GOMAXPROCS(runtime.NumCPU())
	//设置通道数
	channel1 := make(chan bool, 10)
	//开启10个并发
	for i := 0; i < 10; i++ {
		go Go1(channel1, i)
	}
	//读取释放通道，也可理解为等待通道数据
	//若不加等待通道数据，则可能会出现漏执行的数据，即出现异常
	for i := 0; i < 10; i++ {
		<-channel1
	}
}

func Go1(channel1 chan bool, index int) {
	var1 := 1
	for i := 0; i < 1000000; i++ {
		var1 += 1
	}
	f.Println(index, var1)
	//发送通道数据，代表执行完毕
	channel1 <- true
}

//解决并发异常2
func concurrencyFunc2() {
	//设置进程并发数
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := sync.WaitGroup{}
	//开启10个并发
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go Go2(&wg, i)
	}
	wg.Wait()
}

func Go2(wg *sync.WaitGroup, index int) {
	var1 := 1
	for i := 0; i < 1000000; i++ {
		var1 += 1
	}
	f.Println(index, var1)
	wg.Done()
}

//并发select
func concurrencyFunc3() {
	//创建无缓存通道
	channel1, channel2 := make(chan int), make(chan string)
	//创建有缓存通信通道
	channel3 := make(chan bool, 2)
	//执行goroutine
	go func() {
		for {
			select {
			case v, ok := <-channel1:
				if !ok {
					//如果channel被关闭，则在channel3中放入一个true,同时跳出select
					channel3 <- true
					break
				}
				//如果channel未被关闭，则输出语句和通道信息
				f.Println("Channel1", v)
			case v, ok := <-channel2:
				if !ok {
					channel3 <- true
					break
				}
				f.Println("Channel2", v)
			}
		}
	}()
	//向通道传入信息
	channel1 <- 1
	channel2 <- "Hi"
	channel1 <- 3
	channel2 <- "Hello"
	//关闭channel
	close(channel1)
	close(channel2)
	for i := 0; i < 2; i++ {
		//等待通道数据
		<-channel3
	}
}

//并发select2  同时有多个可用的channel时按随机按顺序处理
func concurrencyFunc4() {
	channel := make(chan int)
	go func() {
		for v := range channel {
			f.Println(v)
		}
	}()
	for i := 0; i < 10; i++ {
		select {
		case channel <- 0:
		case channel <- 1:
		}
	}
}

//课后作业，创建一个goroutine，与主线程按顺序相互发送信息若干次，并打印
var channel chan string

func concurrencyFunc5() {
	channel = make(chan string)
	go Go3()
	for i := 0; i < 10; i++ {
		channel <- f.Sprintf("From main：Hello，#%d", i)
		f.Println(<-channel)
	}
}
func Go3() {
	i := 0
	for {
		f.Println(<-channel)
		channel <- f.Sprintf("From main：Hello，#%d", i)
		i++
	}
}

//slice的坑
func sliceFunc() {
	func1 := func(s []int) []int { //最好修改添加容量时设置返回值
		//超出范围地址改变
		s = append(s, 3)
		return s
	}
	s := make([]int, 0)
	f.Println(s)
	s = func1(s)
	//原有地址的值并未改变
	f.Println(s)
}

//time的坑
func timeFunc() {
	t := time.Now()
	f.Println("必须使用它自带的常量的时间值来设置，否则将出现时间偏差。")
	f.Println(t.Format(time.RFC3339))
	f.Println(t.Format("2006-01-02 15:04:05"))
	f.Println(t.Format("2006-01-01 15:04:05"))
}

//闭包的坑
func clouseFunc() {
	s := []string{"a", "b", "c"}
	for _, v := range s {
		//go func() {
		//	//传递的为引用值
		//	f.Println(v)
		//}()
		//使用传参解决
		go func(v string) {
			//传递的为引用值
			f.Println(v)
		}(v)
	}
	select {}
}

//总结与坑
func finishFunc() {
	f.Println("slice的坑：")
	sliceFunc()
	f.Println("time的坑：")
	timeFunc()
	f.Println("闭包的坑：")
	clouseFunc()
}

func getSequence() func() int {
	i := 0
	return func() int {
		i += 1
		return i
	}
}
func otherFunc() {

	const name, age = "Zero", 21
	s := f.Sprintf("%s is %d years old.\n", name, age)
	io.WriteString(os.Stdout, s)

	f.Println("闭包：")
	/* nextNumber 为一个函数，函数 i 为 0 */
	nextNumber := getSequence()

	/* 调用 nextNumber 函数，i 变量自增 1 并返回 */
	f.Println(nextNumber()) //1
	f.Println(nextNumber()) //2
	f.Println(nextNumber()) //3

	/* 创建新的函数 nextNumber1，并查看结果 */
	nextNumber1 := getSequence()
	f.Println(nextNumber1()) //1
	f.Println(nextNumber1()) //2

	f.Println("Hash Map实现：")
	hashMap()
	f.Println("接口实现的类似重载功能：")
	interfaceFunc()
	f.Println("Go的错误处理：")
	_, t := Sqrt(-1)
	f.Println(t)
}

type Phone interface {
	call()
}

type NokiaPhone struct {
}

func (nokiaPhone NokiaPhone) call() {
	f.Println("I am Nokia, I can call you!")
}

type IPhone struct {
}

func (iPhone IPhone) call() {
	f.Println("I am iPhone, I can call you!")
}

func interfaceFunc() {
	var phone Phone

	phone = new(NokiaPhone)
	phone.call()

	phone = new(IPhone)
	phone.call()

}

//go的错误处理
func Sqrt(f float64) (float64, error) {
	if f < 0 {
		return 0, errors.New("math: square root of negative number")
	}
	// 实现
	return 0, errors.New("sdf")
}

func simpleWebFunc() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/login", index)
	http.ListenAndServe("localhost:8080", nil)
}
func index(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	f.Fprint(rw, "Hello World!")
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
	t, _ := template.ParseFiles("huqi/Day_05/index.tpl")
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
