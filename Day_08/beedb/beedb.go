package main

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beedb"
	"github.com/astaxie/beego/orm"
	_ "github.com/ziutek/mymysql/godrv"
	"log"
	"time"
)

/*
beedb是支持database/sql标准接口的ORM库
	Mysql:github.com/ziutek/mymysql/godrv[*]
	Mysql:code.google.com/p/go-mysql-driver[*]
	PostgreSQL:github.com/bmizerany/pq[*]
	SQLite:github.com/mattn/go-sqlite3[*]
	MS ADODB: github.com/mattn/go-adodb[*]
	ODBC: bitbucket.org/miquella/mgodbc[*]

安装
	beedb支持go get方式安装，是完全按照Go Style的方式来实现的。
	go get github.com/astaxie/beedb
*/

//beedb针对驼峰命名会自动帮你转化成下划线字段，例如你定义了Struct名字为UserInfo，那么转化成底层实现的时候是user_info，字段命名也遵循该规则。
type Userinfo struct {
	Uid int `beedb:"PK"` //如果表的主键不是id，那么需要加上pk注释，显式的说这个字段是主键
	//Id         int       `beedb:"PK" sql:"UID" tname:"USER_INFO"` //sql为字段名，tname为定义注释
	Username   string    `orm:"size(5000)" orm:"index" sql:"USERNAME"`
	Departname string    `sql:"DEPARTNAME"`
	Created    time.Time `sql:"CREATED"`
}

//BeeDb初始化连接
func BeedbInit() beedb.Model {
	db, err := sql.Open("mymysql", "gomysql/root/ok") //数据库名、用户名、密码
	if err != nil {
		panic(err)
	}
	//beedb的New函数实际上应该有两个参数，第一个参数标准接口的db，第二个参数是使用的数据库引擎，如果你使用的数据库引擎是MySQL/Sqlite,那么第二个参数都可以省略。
	//如果你使用的数据库是SQLServer，那么初始化需要：
	//orm = beedb.New(db, "mssql")
	//如果你使用了PostgreSQL，那么初始化需要：
	//orm = beedb.New(db, "pg")
	//目前beedb支持打印调试，你可以通过如下的代码实现调试
	//beedb.OnDebug=true
	return beedb.New(db)
}

func InsertDb1(orm beedb.Model) (saveone Userinfo) {
	//var saveone Userinfo
	saveone.Username = "Test Add User"
	saveone.Departname = "Test Add Departname"
	saveone.Created = time.Now()
	rs := orm.Save(&saveone)
	//插入之后saveone.Uid就是插入成功之后的自增ID。Save接口会自动帮你存进去。
	fmt.Println(rs)
	fmt.Println(saveone)
	return
}

func InsertDb2(orm beedb.Model) {
	add := make(map[string]interface{})
	add["username"] = "astaxie"
	add["departname"] = "cloud develop"
	add["created"] = "2012-12-02"
	orm.SetTable("userinfo").Insert(add)
}

//插入多条数据
func InsertDb3(orm beedb.Model) {
	addslice := make([]map[string]interface{}, 0)
	add := make(map[string]interface{})
	add2 := make(map[string]interface{})
	add["username"] = "astaxie"
	add["departname"] = "cloud develop"
	add["created"] = "2012-12-02"
	add2["username"] = "astaxie2"
	add2["departname"] = "cloud develop2"
	add2["created"] = "2012-12-02"
	addslice = append(addslice, add, add2)
	//SetTable函数是显式的告诉ORM，我要执行的这个map对应的数据库表是userinfo
	orm.SetTable("userinfo").InsertBatch(addslice)
}

func UpdateDb1(orm beedb.Model, saveone Userinfo) {
	saveone.Username = "Update Username"
	saveone.Departname = "Update Departname"
	saveone.Created = time.Now()
	orm.Save(&saveone) //现在saveone有了主键值，就执行更新操作
}

func UpdateDb2(orm beedb.Model, saveone Userinfo) {
	t := make(map[string]interface{})
	t["username"] = "astaxie"
	orm.SetTable("userinfo").SetPK("uid").Where(saveone.Uid).Update(t)
	//SetPK：显式的告诉ORM，数据库表userinfo的主键是uid。
	//Where:用来设置条件，支持多个参数，第一个参数如果为整数，相当于调用了Where("主键=?",值)。 Updata函数接收map类型的数据，执行更新数据。
	orm.Save(&saveone) //现在saveone有了主键值，就执行更新操作
}

func SelectDb1(orm beedb.Model, user1 Userinfo) (user2 Userinfo) {
	//Where接受两个参数，支持整形参数
	orm.Where("uid = ?", user1.Uid).Find(&user2)
	//orm.Where(3).Find(&user2) // 这是上面版本的缩写版，可以省略主键，若非主键则必须写出
	//更加复杂的条件：
	//orm.Where("name = ? and age < ?", "john", 88).Find(&user4)
	return
}

//Limit:支持两个参数，第一个参数表示查询的条数，第二个参数表示读取数据的起始位置，默认为0。
//OrderBy:这个函数用来进行查询排序，参数是需要排序的条件。
func SelectDb2(orm beedb.Model, userid int) (allusers []Userinfo) {
	err := orm.Where("id > ?", userid).Limit(10, 20).FindAll(&allusers)
	//省略limit第二个参数，默认从0开始，获取10条数据
	//err := orm.Where("id > ?", "3").Limit(10).FindAll(&tenusers)
	//获取全部数据，条件查询
	//err := orm.OrderBy("uid desc,username asc").FindAll(&everyone)
	//接口函数Select，这个函数用来指定需要查询多少个字段。默认为全部字段*。
	//FindMap()函数返回的是[]map[string][]byte类型，所以你需要自己作类型转换。
	//a, _ := orm.SetTable("userinfo").SetPK("uid").Where(2).Select("uid,username").FindMap()
	log.Fatalln(err)
	return
}

func main() {
	orm := BeedbInit()
	//插入数据操作
	//插入一条记录
	saveone := InsertDb1(orm)
	fmt.Println(saveone)
}

func main2() {
	//初始化连接
	orm := BeedbInit()
	//插入数据操作
	//插入一条记录
	saveone := InsertDb1(orm)
	//map数据插入
	InsertDb2(orm)
	//插入多条数据
	InsertDb3(orm)

	//更新数据
	UpdateDb1(orm, saveone)
	//map操作，更新数据
	UpdateDb2(orm, saveone)

	//查询数据
	//获取一条数据：
	SelectDb1(orm, Userinfo{Uid: 1})
	//获取多条数据，复杂查询
	alluser := SelectDb2(orm, 3)

	//删除数据
	//删除单条数据
	//saveone就是上面示例中的那个saveone
	orm.Delete(&saveone)
	//删除多条数据
	//alluser就是上面定义的获取多条数据的slice
	orm.DeleteAll(&alluser)
	//根据sql删除数据
	orm.SetTable("userinfo").Where("uid > ?", 3).DeleteRow()

	//关联查询
	joinRs, _ := orm.SetTable("userinfo").Join("LEFT", "userdeatail",
		"userinfo.uid=userdeatail.uid").Where("userinfo.uid=?", 1).
		Select("userinfo.uid,userinfo.username,userdeatail.profile").FindMap()
	fmt.Println(joinRs)
	//接口Join函数，这个函数带有三个参数
	//	第一个参数可以是：INNER, LEFT, OUTER, CROSS等
	//	第二个参数表示连接的表
	//	第三个参数表示连接的条件

	//Group By和Having
	groupHavingRs, _ := orm.SetTable("userinfo").GroupBy("username").Having("username='astaxie'").FindMap()
	fmt.Println(groupHavingRs)
	//GroupBy:用来指定进行groupby的字段
	//Having:用来指定having执行的时候的条件

}

func Init() {
	// 需要在init中注册定义的model
	fmt.Println("init db model")

	//orm.RegisterDriver("sqlite3", orm.DRSqlite)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:ok@tcp(127.0.0.1:3306)/gomysql?charset=utf8")
	//orm.RegisterDataBase("default", "sqlite3", "data.db")
	orm.RegisterModel(new(Video), new(Reply))
	orm.RunSyncdb("default", false, true)

}
