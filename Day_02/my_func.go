package main

import (
	f "fmt"
	m "math"
)

const (
	B float64 = 1 << (iota * 10)
	KB
	MB
	GB
	Str1  = 1
	PI    = m.Pi
	Iota1 = iota
	E     = m.E
)

func yunsuanfuFunc() {
	a := 1
	var (
		d = 1
		b = 3
	)

	f.Println(1 << 10)
	f.Println(2048 >> 10)
	f.Println(6 & 11)
	f.Println(6 | 11)
	f.Println(6 ^ 11)
	f.Println(6 &^ 11)
	f.Println(Iota1)      //2 0100
	f.Println(Iota1 << 1) //0010  1248
	f.Println(a, b, d)
	f.Println(B)
	f.Println(KB)
	f.Println(MB)
	f.Println(GB)
}

func pointFunc() {
	a := 1
	a++
	//++a 不可用
	var p *int = &a
	f.Println(p)
	f.Println(*p)
}

func forFunc() {
	a := 1
	//条件语句
	for a < 3 { //等同于white语句
		if a == 1 {
			f.Println("符合条件输出！")
		} else {
			f.Println("不符合条件输出！")
		}
		a++
	}
	for i := 2; i <= 5; i++ {
		f.Println(i)
	}
}

func switchFunc() {

	switch1 := 4
	switch switch1 {
	case 4:
		f.Println("switch1==4")
	default:
		f.Println("默认输出")
	}
	switch { //也可在此处进行初始化变量 switch switch1 := 4{
	case switch1 >= 2:
		f.Println("进入后判断变量值。")
		fallthrough
	case switch1 > 4:
		f.Println("使用fallthrough命令让其判断下一个case。")
	case switch1 <= 5:
		f.Println("上一个为true且未使用fallthrough命令让其判断case。")
	}

}

func labelFun() {
	//跳转语句
LABEL:
	for true {
	Label:
		for i := 0; i < 10; i++ {
			if i == 1 {
				//跳过中间的判断，少输出一个1
				goto goLaebl
			}
			if i > 2 {
				//使用跳转语句跳出死循环
				break LABEL
			} else {
				f.Println(i)
			}
		goLaebl:
			//死循环使用跳转语句跳出
			for true {
				f.Println(i)
				continue Label
			}
		}
	}
}

func arrayFunc() {
	a := [2]int{1, 2}
	f.Println(a)
	b := [2]int{2}
	f.Println("不足长度时自动设置初始值：", b) //[2 0]
	c := [20]int{1, 2, 19: 1}
	f.Println("已知某一下标数时，使用“下标:值”表示：", c) //[1 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1]
	d := [...]int{1, 2, 5: 5, 9: 9}
	f.Println("不确定数组长度时，使用“...”表示：", d) //[1 2 0 0 0 5 0 0 0 9]
	var e *[10]int = &d
	f.Println("指向数组的指针：", e) //&[1 2 0 0 0 5 0 0 0 9]
	g := [...]*int{&a[0], &a[1]}
	f.Println("指向指针的数组：", g) //[0xc0000100e0 0xc0000100e8]
	h1 := [2]int{1, 2}
	h2 := [2]int{1, 3}
	//h3 := [1]int{3}
	f.Println("数组比较：", a == h1)  //true
	f.Println("数组比较：", h1 == h2) //false
	//f.Println("不同类型不可比较：",h2==h3)
	i := new([10]int)
	i[3] = 3
	f.Println("使用new关键字创建指向数组的指针：", i) //&[0 0 0 3 0 0 0 0 0 0]
	j := [2][3]int{{1, 2}, {3, 4}}
	f.Println("多维数组：", j) //[[1 2 0] [3 4 0]]
	k := [...][3]int{{1, 2}, {3, 4}, {5, 6}}
	f.Println("动态多维数组，最内层不可省略指定大小：", k) //[[1 2 0] [3 4 0] [5 6 0]]

}

func maopaoFunc(li [4]int) {
	f.Println(li)
	num := len(li)
	for i := 0; i < num; i++ {
		for j := i + 1; j < num; j++ {
			if li[i] > li[j] {
				tmp := li[i]
				li[i] = li[j]
				li[j] = tmp
			}
		}
	}
	f.Println(li)
}

func sliceFunc() {
	var si []int
	f.Println("直接声明Slice：", si)
	s1 := make([]int, 3, 10)
	f.Println("s1：", s1, "长度：", len(s1), "容量：", cap(s1))
	s2 := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n'}
	sa := s2[2:5]
	f.Println(len(sa), cap(sa))
	sb := sa[9:11]
	f.Println(string(sb))
	s3 := make([]int, 3, 6)
	f.Printf("%p\n", s3)
	s3 = append(s3, 1, 2, 3)
	f.Printf("容量还可存放：%v %p\n", s3, s3)
	s3 = append(s3, 1, 2, 3)
	f.Printf("超出容量开拓新的空间地址：%v %p\n", s3, s3)
	s4 := s2[2:5]
	s5 := s2[1:3]
	f.Println(s4, s5)
	s4[0] = 9
	f.Println("多个切片指向一个相同地址时，改变地址数据，所有切片数据都会改变：", s4, s5)
	s2 = append(s2, 'o', 'p', 'q', 'r', 's', 't')
	s4[1] = 9
	f.Println("但当切片扩容后地址改变，切片数据就不会同步：", s4, s5)
	s6 := []int{1, 2, 3, 4, 5, 6}
	s7 := []int{7, 8, 9}
	copy(s6, s7)
	f.Println("第二个替换到第一个里，若第二个长，则以第一个长度取第二个：", s6) //[7 8 9 4 5 6]
	s8 := s6[:]
	f.Println("完整截取：", s8)
}
