package main

import (
	"time"
)

type SQLModel struct {
	Id       int       `beedb:"PK" sql:"id"`
	Created  time.Time `sql:"created"`
	Modified time.Time `sql:"modified"`
}
type User struct {
	//SQLModel2 SQLModel `beedb:"HasOne"`
	SQLModel `sql:",inline"`
	Name     string `sql:"name" tname:"fn_group"`
	Auth     int    `sql:"auth"`
}

// the SQL table has the columns: id, name, auth, created, modified
// They are marshalled and unmarshalled automatically because of the inline keyword

//ORM主键与外键定义
type Video struct {
	Vid    string `orm:"pk"`
	Name   string
	Replys []*Reply `orm:"reverse(many)"` // 标注反向关系 设置一对一反向关系(可选)
}

type Reply struct {
	Id      int
	Cid     string
	Text    string
	Videoid *Video `orm:"rel(fk)"`
	//Profile *Profile `orm:"rel(one)"`      // OneToOne relation  设置一对多关系
}

//主要是负责数据库处理（mvc中的mode层）
type Store struct {
	Id              int64
	Title           string
	Created         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}

type Customer struct {
	Id              int64
	Uid             int64
	Title           string
	Content         string `orm:"size(5000)"`
	Attachment      string
	Created         time.Time `orm:"index"`
	Updated         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time `orm:"index"`
	ReplyCount      int64
	ReplyLastUserId int64
}

type User2 struct { // 对应user表
	Id      int
	Name    string
	Profile *Profile `orm:"rel(one)"`      // OneToOne relation
	Post    []*Post  `orm:"reverse(many)"` // 设置一对多的反向关系
}
type Profile struct {
	Id   int
	Age  int16
	User *User2 `orm:"reverse(one)"` // 设置一对一反向关系(可选)
}
type Post struct {
	Id    int
	Title string
	User  *User2 `orm:"rel(fk)"` //设置一对多关系
	Tags  []*Tag `orm:"rel(m2m)"`
}
type Tag struct {
	Id    int
	Name  string
	Posts []*Post `orm:"reverse(many)"`
}
