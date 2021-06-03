package main

import (
	f "fmt"
	"sort"
	"strconv"
)

func mapFunc() {
	//初始化
	var map1 map[int]string
	map1 = map[int]string{}
	f.Println("Init Map1：", map1)
	var map2 map[int]string = map[int]string{1: "123", 2: "12312"}
	f.Println("Init Map2：", map2)
	var map3 map[int]string = make(map[int]string)
	f.Println("Init Map3：", map3)
	map4 := make(map[int]string)
	f.Println("Init Map4：", map4)
	//添加元素
	map3[1] = "map3_key1"
	map3[2] = "map3_key2"
	//获取元素
	map3_val := map3[2]
	f.Println("Get Map3_val：", map3_val)
	f.Println("Get Map3：", map3)
	//删除元素
	delete(map3, 2)
	f.Println("Delete Map3_val：", map3)
	f.Println("Get Map3：", map3)
	//复杂Map 内层Map需要逐步独立初始化
	map5 := make(map[int]map[int]string)
	map5[1] = make(map[int]string)
	map5[1][2] = "OK"
	//使用多返回值判断是否初始化，第二个值返回一个bool值
	map5_2_2_val, ok := map5[2][2]
	if !ok {
		f.Println("判断为未初始化，初始化后再赋值！")
		map5[2] = make(map[int]string)
	}
	map5[2][1] = "Good"
	map5_2_2_val = map5[2][1]
	map5_1_2_val := map5[1][2]
	f.Println("提取出复杂数组map5_1_2_val的值：", map5_1_2_val)
	f.Println("提取出复杂数组map5_2_2_val的值：", map5_2_2_val)
	f.Println("Get Map5：", map5)
	li1 := [...]int{1, 5, 8, 4, 5}
	slice1 := li1[:]
	for i, v := range slice1 {
		f.Println("逐个提取切片数据：", v, slice1[i])
	}
	map6 := make([]map[int]string, 5)
	for i := range map6 {
		map6[i] = make(map[int]string, 1)
		map6[i][i] = "map6_" + strconv.Itoa(i)
		f.Println("对Map进行赋值：", map6[i])
	}
	f.Println("Get Map6：", map6)
	//Map排序
	map7 := map[int]string{9: "a", 4: "d", 1: "e", 8: "b", 3: "c"}
	//辅助数组
	li2 := make([]int, len(map7))
	i := 0
	for k, _ := range map7 {
		li2[i] = k
		i++
	}
	sort.Ints(li2)
	f.Println("辅助数组排序后的结果：", li2)
	f.Println("Get Map7：", map7)
	map8 := map[int]string{0: "j", 8: "h", 3: "c", 9: "i", 5: "e", 7: "g", 6: "f", 4: "d", 2: "b", 1: "a"}
	map9 := map[string]int{}
	for k, v := range map8 {
		map9[v] = k
	}
	f.Println("Get Map8：", map8)
	f.Println("Get Map9：", map9)
}

//连续过个变量是相同的类型，则只需定义最后一个的变量的类型即可
func funcFunc1(int1, int2, int3 int) (int4 int, string1, string2 string) {
	f.Println(int1, int2, int3)
	int4, string1, string2 = 1, "2", "3" //返回值定义后，默认方法体内也就不用再定义了，直接赋值即可
	return                               //可以不写，但最好写上return int4, string1, string2
}

/*
	不定长变参，必须放放到参数列表最后，且其为值拷贝，函数内改变不影响原有变量，即形参
	若传入的是切片或数组，为内存地址拷贝，则为实参，改变同样会改变原有值
*/
func funcFunc2(int1 *int, slice1 []int, li1 [4]int, args1 ...int) []int {
	*int1 = 10
	slice1[2] = 10
	li1[2] = 10
	args1[2] = 10
	f.Println("改变接收到的Slice：", slice1)
	f.Println("改变接收到的List：", li1)
	f.Println("改变接收到的Args：", args1)
	return slice1
}

func funcFunc3(x int) func(int) int {
	return func(y int) int {
		f.Println("返回闭包函数！")
		return x * y
	}
}

func funcFunc4() {
	f.Println("defer函数，逆序执行：")
	for i := 0; i < 5; i++ {
		defer f.Println(i)
		defer func() {
			f.Println("匿名闭包函数，i值为地址引用，i都为相同的最终结果：", i)
		}()
	}
}

func structFunc() {
	type Person1 struct {
		Name string
		Age  int
	}
	person1 := &Person1{}
	person2 := &Person1{Name: "Wonder", Age: 20}
	func1 := func(person3 *Person1) {
		person3.Age = 21
		f.Println(person3)
	}
	func1(person2)
	person1.Name = "Zero"
	person1.Age = 19
	f.Println(person1, person2)

	f.Println("Struct匿名结构、嵌套结构与匿名字段：")
	person4 := &struct {
		Id, Age int
		Name    string
		Deatile struct {
			Phone, City string
		}
	}{Id: 1, Age: 20, Name: "Struct1"}
	person4.Deatile.Phone = "18888888888"
	person4.Deatile.City = "江苏省南京市建邺区"
	person5 := &struct {
		int
		string
	}{5, "Person5"}
	f.Println(person4, person5)
	type Person2 struct {
		Name string
		Age  int
	}
	person6 := Person2{Name: "Wonder", Age: 20}
	person7 := Person2{Name: "Wonder", Age: 20}
	//结构体内容虽相同，但名字不同，即不可比较
	//f.Println(person2 == person6)
	f.Println(person7 == person6) //true
	f.Println("")
	type Student struct {
		Person2
		Role int
		Age  int
	}
	student1 := Student{Person2: Person2{Age: 21, Name: "Zero"}, Role: 1}
	student1.Role = 1
	student1.Name = "Wonder"
	//只要同级结构体不重名就不会出错
	student1.Age = 19
	student1.Person2.Age = 20
	f.Println("结构体嵌套和重名时：", student1)
}

type Person struct {
	Name string
	Age  int
}
type SI int

//相当于给结构体加入方法，方法与结构体绑定，不可进行重载
func (person *Person) isTrueAge() bool {
	f.Println(person.Name)
	return person.Age > 1
}
func (si *SI) Increase(num int) {
	*si += SI(num)
	//*si = *si + 100
}
func methodFunc() {
	person1 := Person{Name: "Zero", Age: 21}
	f.Println("Method Value：", person1.isTrueAge())
	f.Println("Method Expression：", (*Person).isTrueAge(&person1))

	//在同一个包中，所有类型和变量都为公有。只有当不处于同一个包中时，才有访问权限问题，首字母大写才可访问。
	f.Println("课堂作业：")
	var si SI = 0
	si.Increase(100)
	f.Println(si)
}

type USB interface {
	Name() string
	Connecter
}
type Connecter interface {
	Connect()
}
type PhoneConnecter struct {
	name string
}

//接口中的方法必须被全部实现
func (pc PhoneConnecter) Name() string {
	return pc.name
}
func (pc PhoneConnecter) Connect() {
	f.Println("Connect：", pc.name)
}

//func Disconnect(usb USB)  {
//	if pc, ok := usb.(PhoneConnecter); ok{
//		f.Println("Disconnected：", pc.name)
//		return
//	}
//	f.Println("Unknown decive.")
//}
func Disconnect(usb interface{}) {
	//使用接收广泛类型，调用后在进行具体判断
	switch v := usb.(type) {
	case PhoneConnecter:
		f.Println("Disconnected：", v.name)
	default:
		f.Println("Unknown decive.")
	}
}
func interfaceFunc() {
	var usb USB
	usb = PhoneConnecter{"PhoneConnecter"}
	usb.Connect()
	f.Println(usb.Name())
	Disconnect(usb)
	pc := PhoneConnecter{"ComputerConnecter"}
	var con Connecter
	//强制类型转换，大的可以转换成小的，即降级转换
	con = Connecter(pc)
	con.Connect()
	//只有当接口存储的类型和对象都为nil时，接囗才等于nil
	var intI interface{}
	f.Println(intI == nil)
	var p *int = nil
	intI = p
	f.Println(intI == nil)
}
