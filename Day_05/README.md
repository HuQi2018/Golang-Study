# Go语言学习笔记 Day05

### [Day04](http://njgit.jsaepay.com/wangwei/golang-study/src/branch/master/huqi/Day_04)
### [Day06](http://njgit.jsaepay.com/wangwei/golang-study/src/branch/master/huqi/Day_06)


### 1. 并发Concurrency
	Go语言具有高并发
	Goroutine只是由官方实现的超级“线程池”而已。
	每个实例4-5KB的栈内存占用和由于实现机制而大幅减少的创建和销毁开销，是制造Go号称高并发的根本原因
	并发主要由切换时间片来实现“同时”运行，在并行则是直接利用多核实现多线程的运行，但Go可以设置使用核数，以发挥多核计算机的能力。
	Goroutinue奉行通过通信来共享内存，而不是共享内存来通信。

	Channel
		Channel是goroutine沟通的桥梁，大都是阻塞同步的
		通过make创建，close关闭
		Channel是引用类型
		可以使用for range来迭代不断操作channel
		可以设置单向或双向通道
		可以设置缓存大小，在未被填满前不会发生阻塞

	Select
		可以处理一个或多个channel的发送与接收
		同时有多个可用的channel时按随机按顺序处理
		可用空的select来阻塞main函数
		可设置超时

	channel分为两种：一种是有buffer的，一种是没有buffer的，默认是没有buffer的
	ci := make(chan int) //无buffer
	cj := make(chan int, 0) //无buffer
	cs := make(chan int, 100) //有buffer
	有缓冲的channel，因此要注意放入数据的操作c<- 0先于取数据操作 <-c
	无缓冲的channel，因此要必须保证取操作<-c 先于放操作c<- 0
[Go语言_并发篇](https://www.cnblogs.com/yjf512/archive/2012/06/06/2537712.html)

```go
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
		<- channel1
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

func Go2(wg *sync.WaitGroup, index int)  {
	var1 := 1
	for i := 0; i < 1000000; i++ {
		var1 += 1
	}
	f.Println(index, var1)
	wg.Done()
}
//并发select
func concurrencyFunc3()  {
	//创建无缓存通道
	channel1, channel2 := make(chan int), make(chan string)
	//创建有缓存通信通道
	channel3 := make(chan bool, 2)
	//执行goroutine
	go func() {
		for {
			select {
			case v, ok := <- channel1:
				if !ok {
					//如果channel被关闭，则在channel3中放入一个true,同时跳出select
					channel3 <- true
					break
				}
				//如果channel未被关闭，则输出语句和通道信息
				f.Println("Channel1", v)
			case v, ok := <- channel2:
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
		<- channel3
	}
}
//并发select2  同时有多个可用的channel时按随机按顺序处理
func concurrencyFunc4()  {
	channel := make(chan int)
	go func() {
		for v := range channel{
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
func concurrencyFunc5()  {
	channel = make(chan string)
	go Go3()
	for i := 0; i < 10; i++ {
		channel <- f.Sprintf("From main：Hello，#%d", i)
		f.Println(<-channel)
	}
}
func Go3()  {
	i := 0
	for {
		f.Println(<-channel)
		channel <- f.Sprintf("From main：Hello，#%d", i)
		i++
	}
}
```

### 2. 总结与其他补充
	GOPATH与GOROOT路径不能相同
	GOROOT指向安装路径
	GOPATH中有bin、pkg、src（必有）文件夹

	select 语句	select 语句类似于 switch 语句，但是select会随机执行一个可运行的case。如果没有case可运行，它将阻塞，直到有case可运行。
	注意：Go 没有三目运算符，所以不支持 ?: 形式的条件判断。
	range 关键字用于 for 循环中迭代数组(array)、切片(slice)、通道(channel)或集合(map)的元素。
	Map 是无序的，我们无法决定它的返回顺序，这是因为 Map 是使用 hash 表来实现的。
	/* 声明变量，默认 map 是 nil */
	var map_variable map[key_data_type]value_data_type
	/* 使用 make 函数 */
	map_variable := make(map[key_data_type]value_data_type)
	如果不初始化 map，那么就会创建一个 nil map。nil map 不能用来存放键值对。
	go 不支持隐式转换类型
	
	Go 并发：
	同一个程序中的所有 goroutine 共享同一个地址空间。
	操作符 <- 用于指定通道的方向，发送或接收。如果未指定方向，则为双向通道。

```go
//slice的坑
func sliceFunc() {
	func1 := func(s []int) []int{ //最好修改添加容量时设置返回值
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
	for _, v := range s{
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
	select {

	}
}

//Go基础语言的总结与坑
func finishFunc()  {
	f.Println("slice的坑：")
	sliceFunc()
	f.Println("time的坑：")
	timeFunc()
	f.Println("闭包的坑：")
	clouseFunc()
}

//go的类似重载
type Phone interface {
    call()
}

type NokiaPhone struct {
}

func (nokiaPhone NokiaPhone) call() {
    fmt.Println("I am Nokia, I can call you!")
}

type IPhone struct {
}

func (iPhone IPhone) call() {
    fmt.Println("I am iPhone, I can call you!")
}

func main() {
    var phone Phone

    phone = new(NokiaPhone)
    phone.call()

    phone = new(IPhone)
    phone.call()

}

```

###[Go的目录结构标准](https://github.com/golang-standards/project-layout/blob/master/README_zh-CN.md)
###[Go的相关项目地址](https://github.com/golang/go/wiki/Projects)
###[etcd Go的键值存储数据库项目](https://github.com/etcd-io/etcd)
###[nsq Go的消息传递项目](https://github.com/nsqio/nsq)
###[beego Web框架](https://beego.me/)
###[revel Web框架](http://revel.github.io/)


### 3. 简单的网络服务
```go
func simpleWebFunc() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/login", index)
	http.ListenAndServe("localhost:8080", nil)
}
func index(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	f.Fprint(rw, "Hello World!")
	message := ""
	if req.Method == "POST"{
		req.ParseForm()
		if cheackUserName(req.Form["username"][0]){
			f.Fprint(rw, "用户名：", req.Form["username"][0], "密码：", req.Form["password"][0])
			return
		}else {
			message = "用户名长度不符合要求！"
		}
	}
	t, _ := template.ParseFiles("huqi/Day_05/index.tpl")
	t.Execute(rw, message)
}
func cheackUserName(username string) bool {
	if len(username) >= 6 && len(username) <= 16{
		return true
	}
	return false
}
func handler(rw http.ResponseWriter, req *http.Request) {
	//解析用户请求的数据
	req.ParseForm()
	//判断用户是否有请求参数
	if len(req.Form["name"]) > 0{
		f.Fprint(rw, "Hello ", req.Form["name"][0])
	}else {
		f.Fprint(rw, "Hello World!")
	}
}
```
