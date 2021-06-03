# Go语言学习笔记 Day03

### [Day02](http://njgit.jsaepay.com/wangwei/golang-study/src/branch/master/huqi/Day_02)
### [Day04](http://njgit.jsaepay.com/wangwei/golang-study/src/branch/master/huqi/Day_04)

### 1. Map类型

	类似其它语言中的哈希表或者字典，以key-value形式存储数据
	Key必须是支持==或!=比较运算的类型，不可以是函数、map或slice
	Map查找比线性搜索快很多，但比使用索引访问数据的类型慢100倍
	Map使用make()创建，支持:=这种简写方式
	make([keyType]valueType,cap)，cap表示容量，可省略
	超出容量时会自动扩容，但尽量提供一个合理的初始值
	使用len()获取元素个数
	键值对不存在时自动添加，使用delete()删除某踺值对
	使用for range对map和slice进行迭代操作

```go
	//初始化
	var map1 map[int]string
	map1 = map[int]string{}
	f.Println("Init Map1：", map1)
	var map2 map[int]string = map[int]string{1:"123",2:"12312"}
	f.Println("Init Map2：", map2)
	var map3 map[int]string = make(map[int]string)
	f.Println("Init Map3：", map3)
	map4 := make(map[int]string)
	f.Println("Init Map4：", map4)
	//添加元素
	map3[1]="map3_key1"
	map3[2]="map3_key2"
	//获取元素
	map3_val := map3[2]
	f.Println("Get Map3_val：" ,map3_val)
	f.Println("Get Map3：" ,map3)
	//删除元素
	delete(map3, 2)
	f.Println("Delete Map3_val：" ,map3)
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
	f.Println("提取出复杂数组map5_2_2_val的值：",map5_2_2_val)
	f.Println("Get Map5：",map5)
	li1 := [...]int{1,5,8,4,5}
	slice1 := li1[:]
 	for i,v := range slice1{
		f.Println("逐个提取切片数据：", v,slice1[i])
	}
	map6 := make([]map[int]string, 5)
	for i := range map6 {
		map6[i] = make(map[int]string, 1)
		map6[i][i] = "map6_" + strconv.Itoa(i)
		f.Println("对Map进行赋值：", map6[i])
	}
	f.Println("Get Map6：", map6)
	//Map排序
	map7 := map[int]string{9:"a", 4:"d", 1:"e", 8:"b", 3:"c"}
	//辅助数组
	li2 := make([]int, len(map7))
	i := 0
	for k, _ := range map7{
		li2[i] = k
		i++
	}
	sort.Ints(li2)
	f.Println("辅助数组排序后的结果：", li2)
	f.Println("Get Map7：", map7)
	map8 := map[int]string{0:"j", 8:"h", 3:"c", 9:"i", 5:"e", 7:"g", 6:"f", 4:"d", 2:"b", 1:"a"}
	map9 := map[string]int{}
	for k, v := range map8{
		map9[v] = k
	}
	f.Println("Get Map8：", map8)
	f.Println("Get Map9：", map9)
```


### 2. Go函数function

	Go函数不支持嵌套、重载和默认参数
	但支持一下特性：
		无需声明原型、不定长度变参、多返回值、命名返回值参数、匿名函数、闭包
	定义函数使用关键字func，且左大括号不能另起一行
	函数也可以作为一种类型使用

	defer
		defer的执行方式类似于其他语言中的析构函数，在函数体执行结束后按照调用顺序的相反顺序逐个执行
		即使函数发生严重错误也会执行
		支持匿名函数的调用
		常用于资源清理、文件关闭、解锁以及记录时间等操作
		通过与匿名函数配合可在return之后修改函数的计算结果
		如果函数体内某个变量作为defer时匿名函数的参数，则在定义defer时即已经获得了拷贝，否则则是引用某个变量的地址
		Go没有异常机制，但有panic/recover模式来处理错误
		panic可以在任何地方引发，但recover只有在defer调用的函数中有效

```go

//func.go

//连续过个变量是相同的类型，则只需定义最后一个的变量的类型即可
func funcFunc1(int1, int2, int3 int) (int4 int, string1, string2 string) {
	f.Println(int1, int2, int3)
	int4, string1, string2 = 1, "2", "3"  //返回值定义后，默认方法体内也就不用再定义了，直接赋值即可
	return //可以不写，但最好写上return int4, string1, string2
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

//main.go

f.Println("Func函数练习1：")
	f.Println(funcFunc1(1, 2, 3))
	f.Println("")

	f.Println("Func函数练习2：")
	slice1 := []int{1,2,3,4}
	li1 := [4]int{5,6,7,8}
	int1 := 1
	func1 := funcFunc2(&int1,slice1, li1, 4, 5, 6)
	f.Println("接收到的返回值：", func1)
	f.Println("Int传地址指针会改变：", int1)
	f.Println("Slice会改变：", slice1)
	f.Println("List不会改变：", li1)
	func2 := func(int2 int) {
		f.Println("输出匿名函数！")
	}
	func2(2)
	func3 := funcFunc3(5)
	f.Println(func3(5))
	f.Println(func3(9))
	f.Println("类似的嵌套函数：")
	nested := func() {
		f.Println("外层函数！")
		deeplyNested := func() {
			f.Println("内层函数！")
		}
		deeplyNested()
	}
	nested()
	f.Println("")
	f.Println("defer函数：")
	//defer函数
	funcFunc4()
	//使用defer处理go语言的错误机制
	func4 := func() {f.Println("Func4！")}
	func5 := func() {
		defer func() {  //将程序从Panic恢复到Recover状态
			if err := recover();err != nil{
				f.Println("Recover in Func5！")
			}
		}()
		panic("Panic in Func5！")  //终止
	}
	func6 := func() {f.Println("Func6！")}
	func4()
	func5()
	func6()
	f.Println("")

	//课后习题
	f.Println("课后习题：")
	fs := [4]func(){}
	for i := 0; i <4; i++ {
		defer f.Println("defer i = ", i)  //使用i的值
		defer func() {f.Println("defer_closure i = ", i)}()  //使用地址
		fs[i] = func() {
			f.Println("closure i = ", i)  //使用地址
		}
	}
	for _, f := range fs{
		f()
	}
	f.Println("")
```


### 3. 结构struct
	Go中的struct与C中的struct非常相似，并且Go没有class
	使用type<Name>struct{}定义结构，名称遵循可见性规则
	支持指向自身的指针类型成员
	支持匿名结构，可用作成员或定义成员变量
	匿名结构也可以用于map的值
	可以使用字面值对结构进行初始化
	允许直接通过指针来读写结构成员
	相同类型的成员可进行直接拷贝赋值
	支持==与!=比较运算符，但不支持>或<
	支持匿名字段，本质上是定义了以某个类型名为名称的字段
	嵌入结构作为匿名字段看起来像继承，但不是继承
	可以使用匿名字段指针

```go
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
	person1.Age  = 19
	f.Println(person1, person2)

	f.Println("Struct匿名结构、嵌套结构与匿名字段：")
	person4 := &struct {
		Id,Age int
		Name   string
		Deatile struct{
			Phone, City string
		}
	}{Id: 1, Age: 20,Name: "Struct1"}
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
	f.Println(person7 == person6)  //true
	f.Println("")
	type Student struct {
		Person2
		Role int
		Age  int
	}
	student1 := Student{Person2: Person2{Age: 21, Name: "Zero"}, Role:1}
	student1.Role = 1
	student1.Name = "Wonder"
	//只要同级结构体不重名就不会出错
	student1.Age = 19
	student1.Person2.Age = 20
	f.Println("结构体嵌套和重名时：", student1)
```

### 4. 方法Method
	Go中虽没有class，但依旧有method
	通过显示说明receiver来实现与某个类型的组合
	只能为同一个包中的类型定义方法
	Receiver可以是类型的值或者指针
	不存在方法重载
	可以使用值或指针来调用方法，编译器会自动完成转换
	从某种意义上来说，方法是函数的语法糖，因为receiver其实就是方法所接收的第一个参数
	如果外部结构和嵌入结构存在同名方法，则优先调用外部结构的方法
	类型别名不会拥有底层类型所附带的方法
	方法可以调用结构中的非公开字段

```go
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
func methodFunc()  {
	person1 := Person{Name: "Zero", Age: 21}
	f.Println("Method Value：", person1.isTrueAge())
	f.Println("Method Expression：", (*Person).isTrueAge(&person1))

	//在同一个包中，所有类型和变量都为公有。只有当不处于同一个包中时，才有访问权限问题，首字母大写才可访问。
	f.Println("课堂作业：")
	var si SI = 0
	si.Increase(100)
	f.Println(si)
}
```

### 5. 接口Interface
	接囗是一个或多个方法签名的集合
	只要某个类型拥有该接囗的所有方法签名即算实现该接囗，无需显示声明实现了哪个接囗，这称为Structural Typing
	接囗只有方法声明，没有实现，没有数据字段
	接囗可以匿名嵌入其它接囗，或嵌入到结构中
	将对象赋值给接囗时，会发生拷贝，而接囗内部存储的是指向这个复制品的指针，既无法修改复制品的状态，也无法获取指针
	只有当接口存储的类型和对象都为nil时，接囗才等于nil
	接囗调用不会做receiver的自动转换
	接囗同样支持匿名字段方法
	接口也可实现类似OOP中的多态
	空接囗可以作为任何类型数据的容器

	类型断言
		通过类型断言的ok pattern可以判断接口中的数据类型
		使用type switch则可针对空接口进行比较全面的类型判断
	接口转换
		可以将拥有超集的接口转换为子集的接口

```go
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
func (pc PhoneConnecter) Name() string{
	return pc.name
}
func (pc PhoneConnecter) Connect(){
	f.Println("Connect：", pc.name)
}
//func Disconnect(usb USB)  {
//	if pc, ok := usb.(PhoneConnecter); ok{
//		f.Println("Disconnected：", pc.name)
//		return
//	}
//	f.Println("Unknown decive.")
//}
func Disconnect(usb interface{})  {
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
```
