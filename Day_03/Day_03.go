package main

import f "fmt"

func main() {
	//f.Println("Map函数练习：")
	//mapFunc()
	//f.Println("")
	//f.Println("Func函数练习1：")
	//f.Println(funcFunc1(1, 2, 3))
	//f.Println("")
	//
	//f.Println("Func函数练习2：")
	//slice1 := []int{1,2,3,4}
	//li1 := [4]int{5,6,7,8}
	//int1 := 1
	//func1 := funcFunc2(&int1,slice1, li1, 4, 5, 6)
	//f.Println("接收到的返回值：", func1)
	//f.Println("Int传地址指针会改变：", int1)
	//f.Println("Slice会改变：", slice1)
	//f.Println("List不会改变：", li1)
	//func2 := func(int2 int) {
	//	f.Println("输出匿名函数！")
	//}
	//func2(2)
	//func3 := funcFunc3(5)
	//f.Println(func3(5))
	//f.Println(func3(9))
	//f.Println("类似的嵌套函数：")
	//nested := func() {
	//	f.Println("外层函数！")
	//	deeplyNested := func() {
	//		f.Println("内层函数！")
	//	}
	//	deeplyNested()
	//}
	//nested()
	//f.Println("")
	//f.Println("defer函数：")
	////defer函数
	//funcFunc4()
	////使用defer处理go语言的错误机制
	//func4 := func() {f.Println("Func4！")}
	//func5 := func() {
	//	defer func() {  //将程序从Panic恢复到Recover状态
	//		if err := recover();err != nil{
	//			f.Println("Recover in Func5！")
	//		}
	//	}()
	//	panic("Panic in Func5！")  //终止
	//}
	//func6 := func() {f.Println("Func6！")}
	//func4()
	//func5()
	//func6()
	//f.Println("")
	//
	////课后习题
	//f.Println("课后习题：")
	//fs := [4]func(){}
	//for i := 0; i <4; i++ {
	//	defer f.Println("defer i = ", i)  //使用i的值
	//	defer func() {f.Println("defer_closure i = ", i)}()  //使用地址
	//	fs[i] = func() {
	//		f.Println("closure i = ", i)  //使用地址
	//	}
	//}
	//for _, f := range fs{
	//	f()
	//}
	//f.Println("")

	//f.Println("结构Struct：")
	//structFunc()
	//f.Println("")

	//f.Println("方法Method：")
	//methodFunc()
	//f.Println("")

	f.Println("接口Interface：")
	interfaceFunc()
	f.Println("")

}
