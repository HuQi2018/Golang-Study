package main

import (
	f "fmt"
	"reflect"
)

type User struct {
	Id   int
	Name string
	Age  int
}

type Manager struct {
	User
	Title string
	Role  string
}

func (u User) Hello(name string) {
	f.Println("Hello ", name, ", my name is", u.Name)
}

func (manger Manager) Admin(title string) string {
	manger.Title = title
	f.Println("接口接收到title参数：", title)
	return manger.Role
}

func reflectInfo1(inter1 interface{}) {
	reflectVar1 := reflect.TypeOf(inter1)
	f.Println("获取反射接口的名字：", reflectVar1.Name())
	//判断传入的类是否正确
	if k := reflectVar1.Kind(); k != reflect.Struct {
		f.Println("传入类型错误！")
		return
	}
	reflectVar2 := reflect.ValueOf(inter1)
	f.Println("获取反射接口的值：")
	for i := 0; i < reflectVar1.NumField(); i++ {
		reflectVar3_i := reflectVar1.Field(i)
		reflectVar3_val := reflectVar2.Field(i).Interface()
		f.Printf("%6s：%v = %v\n", reflectVar3_i.Name, reflectVar3_i.Type, reflectVar3_val)
	}
	for i := 0; i < reflectVar1.NumMethod(); i++ {
		reflectVar4_m := reflectVar1.Method(i)
		f.Printf("%6s：%v\n", reflectVar4_m.Name, reflectVar4_m.Type)
	}
}

func reflectInfo2(inter1 interface{}) {
	reflectVar1 := reflect.TypeOf(inter1)
	f.Printf("%#v\n", reflectVar1.FieldByIndex([]int{0, 1})) //传入切片，得到指定字段的值，将要取到的字段设置为1
}

func reflectInfo3(inter1 interface{}) {
	reflectVar1 := reflect.ValueOf(inter1)
	//判断传入的类是否正确
	if reflectVar1.Kind() != reflect.Ptr || !reflectVar1.Elem().CanSet() {
		f.Println("传入类型错误！")
		return
	} else {
		reflectVar1 = reflectVar1.Elem()
	}
	reflectVar2 := reflectVar1.FieldByName("Name")
	if !reflectVar2.IsValid() {
		f.Println("不存在该字段类型，修改失败！")
		return
	}
	if reflectVar2.Kind() == reflect.String {
		reflectVar2.SetString("反射修改的名字")
	}
}
func reflectInfo4(user User) {
	//反射的结构
	reflectVal1 := reflect.ValueOf(user)
	//反射接口的函数名
	reflectVal2 := reflectVal1.MethodByName("Hello")
	//设置反射函数参数
	args := []reflect.Value{reflect.ValueOf("joe")}
	//调用反射函数
	reflectVal2.Call(args)
}

//打印反射的接口类型
func reflectPrint(intfer interface{}) {
	type1 := reflect.TypeOf(intfer)
	value := reflect.ValueOf(intfer)
	if k := type1.Kind(); k != reflect.Struct {
		f.Println("传入类型错误！")
		return
	}
	for i := 0; i < type1.NumField(); i++ {
		reflectVar3_i := type1.Field(i)
		reflectVar3_val := value.Field(i).Interface()
		type2 := reflect.TypeOf(reflectVar3_val)
		if type2.Kind() == reflect.Struct {
			reflectPrint(reflectVar3_val)
		} else {
			f.Printf("%6s：%v = %v\n", reflectVar3_i.Name, reflectVar3_i.Type, reflectVar3_val)
		}
	}
}

func reflectInfo5(manager interface{}, manager2 interface{}, name string) {
	//反射的结构
	reflectVar1 := reflect.ValueOf(manager)
	//反射接口的函数名
	reflectVar2 := reflectVar1.MethodByName("Admin")
	//设置反射函数参数
	args := []reflect.Value{reflect.ValueOf("title")}
	//调用反射函数
	f.Println(reflectVar2.Call(args))
	//判断传入的类是否正确
	if reflectVar1.Kind() != reflect.Ptr || !reflectVar1.Elem().CanSet() {
		f.Println("传入类型错误！")
		return
	} else {
		f.Println("获取反射接口的值：")
		reflectPrint(manager2)
		reflectVar1 = reflectVar1.Elem()
	}
	reflectVar3 := reflectVar1.FieldByName("Name")
	if !reflectVar3.IsValid() {
		f.Println("不存在该字段类型，修改失败！")
		return
	}
	if reflectVar3.Kind() == reflect.String {
		reflectVar3.SetString(name)
	}
}
func reflectionFunc() {
	f.Println("反射获取类名及值：")
	user1 := User{1, "Zero", 21}
	//传入的值必须为拷贝的值，不能是指针类型
	reflectInfo1(user1)
	f.Println("")
	f.Println("反射获取嵌套类名及值：")
	user2 := Manager{Title: "管理员", User: user1}
	reflectInfo2(user2)
	f.Println("")
	f.Println("使用反射修改变量值：")
	//reflectVal1 := 123
	//reflectVal2 := reflect.ValueOf(&reflectVal1)
	//reflectVal2.Elem().SetInt(999)
	//f.Println(reflectVal1)
	reflectInfo3(&user1)
	f.Println("反射修改后的值：", user1)
	f.Println("")
	f.Println("使用反射进行动态调用：")
	user3 := User{1, "Zero", 21}
	reflectInfo4(user3)
	f.Println("课后练习：")
	user4 := Manager{Role: "Admin", Title: "Manager", User: user3}
	f.Println("反射调用初始化之前：", user4)
	reflectInfo5(&user4, user4, "Wonder")
	f.Println("反射调用之后：", user4)
	f.Println("")
}

func Go() {
	f.Println("Go并发函数！")
}

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
	channel2 := make(chan bool, 1) //设置缓存为1，有缓存则为同步阻塞，无缓存则为异步的
	go func() {
		f.Println("有缓存读！")
		<-channel2
	}()
	channel2 <- true
}
