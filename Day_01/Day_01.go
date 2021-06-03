//当前程序的包名
package main

//import "fmt"
import (
	f "fmt"
	m "math"
	"strconv"
)

//常量的定义
const (
	PI     = 3.14
	const1 = "1"
	const2 = byte(2)
	const3 = m.Pi

	upV1
	upV2

	upA1, upA2 = 1, "2"
	//常量赋值时，下行不赋值时，上下两行数目必须一致，否则编译不通过
	upA3, upA4

	iota1 = iota // 8
)

const (
	iota2 = 'A'  // 65
	iota3 = iota // 1
	iota4 = 'B'  // 66
	iota5 = iota //2
)

//全局变量的声明与赋值
var (
	name  = "gopher"
	name1 = m.MaxInt16
	name2 = 2
	name3 = m.Abs(-3)
	//默认赋值
	default1 int
	default2 float32
	default3 string
	default4 bool
	default5 []int
	default6 []byte
	//使用并行方式以及类型推断
	var1, var2 = 1, m.E
	//不可以省略var
	//var3 := false
)

//一般类型的声明
type (
	newType int
	type1   float32
	type2   string
	type3   byte
	文本      string
)

//结构的声明
type (
	gopher struct {
	}
)

//接口的声明
type (
	golang interface {
	}
)

//由main函数作为程序入口点启动
func main() {
	//使用变量别名
	var otherName 文本 = "中文文本"
	//变量赋值
	var a, b, c, d int
	a, b, c, d = 1, 2, 3, 4
	i, _, n, o := 13, 14, 15, 16
	//变量类型转换
	//在相互兼容的两种类型之间进行转换
	var e float32 = 1.1
	g := int(e)
	//以下表达式无法通过编译
	//var h bool = true
	//j := int(h)
	var h int = 65
	j := string(h)         //A
	k := strconv.Itoa(h)   // 65
	h, _ = strconv.Atoi(k) // 65

	//fmt.Println("Hello World!你好,世界！")
	f.Println("输出引用变量：", m.Pi)
	f.Println("输出变量(可用于判断是否超出变量类型范围)：", name1)
	f.Println("输出变量：", name3)
	f.Println("输出变量：", name2)
	f.Println("输出测试默认值：", default1)
	f.Println("输出测试默认值：", default2)
	f.Println("输出测试默认值：", default3)
	f.Println("输出测试默认值：", default4)
	f.Println("输出测试默认值：", default5)
	f.Println("输出测试默认值：", default6)
	f.Println("使用变量别名：", otherName)
	f.Println("使用并行方式赋值：", var1)
	f.Println("使用并行方式赋值：", var2)
	f.Println("使用局部变量：", a, b, c, d, i, n, o)
	f.Println("变量类型强制转换：", g, j, k, h)
	f.Println("当常量未被赋值时将使用上行表达式的值：", upV1, upV2)
	f.Println("iota的值：", iota1, iota2, iota3, iota4, iota5)
	f.Println("Hello Go!你好,Go！")

}
