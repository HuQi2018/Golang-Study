# Go语言学习笔记 Day07

### [Day06](http://njgit.jsaepay.com/wangwei/golang-study/src/branch/master/huqi/Day_06)
### [Day08](http://njgit.jsaepay.com/wangwei/golang-study/src/branch/master/huqi/Day_08)

------------

####[Go 语言 Web 应用开发](https://github.com/unknwon/building-web-applications-in-go)

####[Go Web](https://github.com/unknwon/building-web-applications-in-go)

####[Go 名库](https://github.com/unknwon/go-rock-libraries-showcases)

### 1. 文本模板引擎
	1、将模板应用于给定的数据结构来执行模板，模板的编码与 Go 语言源代码文件相同，需为 UTF-8 编码
	2、模板中的注解（Annotation）会根据数据结构中的元素来执行并派生具体的显示结构，这些元素一般指结构体中的字段或 map 中的键名
	3、模板的执行逻辑会依据点（Dot，"."）操作符来设定当前的执行位置，并按序完成所有逻辑的执行。
	4、模板中的行为（Action）包括数据评估（Data Evaluation）和控制逻辑，且需要使用双层大括号（{{ 和 }}）包裹。除行为以外的任何内容都会原样输出不做修改。
	5、模板解析完成后，从设计上可以并发地进行渲染，但要注意被渲染对象的并发安全性。例如，一个模板可以同时为多个客户端的响应进行渲染，因为输出对象（Writer）是相互独立的，但是被渲染的对象可能有各自的状态和时效性。

### 今日学习模板处理，具体内容参考代码内容
