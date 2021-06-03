/*
创建者：     Zero
创建时间：   2021/5/26
项目名称：   goWeb
*/
package sys

import "time"

//电影信息
type MovieInfo struct {
	Id            int64
	MovieId       int64     `xorm:"unique index" json:"movie_id"`                //豆瓣电影唯一标识
	Title         string    `xorm:"varchar(1000) notnull index" json:"name"`     //电影名称
	Types         string    `xorm:"varchar(20) notnull" json:"types"`            //电影类型
	Rating        string    `xorm:"longtext notnull" json:"rating"`              //电影评分
	Durations     string    `xorm:"varchar(80) notnull" json:"durations"`        //时长
	Year          int64     `xorm:"varchar(20) notnull" json:"year"`             //上映年份
	Pubdates      string    `xorm:"varchar(1000) notnull" json:"pubdates"`       //上映日期数据
	Pubdate       string    `xorm:"varchar(1000) notnull" json:"pubdate"`        //上映日期
	Aka           string    `xorm:"varchar(1000) notnull" json:"aka"`            //又名
	Tags          string    `xorm:"varchar(1000) notnull" json:"tags"`           //标签
	OriginalTitle string    `xorm:"varchar(1000) notnull" json:"original_title"` //原始标题
	Language      string    `xorm:"varchar(80) notnull" json:"language"`         //电影语言
	Country       string    `xorm:"varchar(1000) notnull" json:"country"`        //制片国家
	Actors        string    `xorm:"longtext notnull" json:"actors"`              //演员
	Writers       string    `xorm:"longtext notnull" json:"writers"`             //作者
	Directors     string    `xorm:"longtext notnull" json:"directors"`           //导演
	Summary       string    `xorm:"longtext notnull" json:"summary"`             //简介
	Photos        string    `xorm:"longtext notnull" json:"photos"`              //图像数据
	Images        string    `xorm:"longtext notnull" json:"images"`              //海报数据
	Videos        string    `xorm:"longtext notnull" json:"videos"`              //视频数据
	OpUser        string    `xorm:"index" json:"op_user"`                        //操作用户人
	CreatedAt     time.Time `xorm:"created" json:"created_at"`
	UpdatedAt     time.Time `xorm:"updated" json:"updated_at"`
	DeletedAt     time.Time `xorm:"deleted" json:"deleted_at"`
}

//电影类型
type MovieType struct {
	Id        int64
	Name      string    `xorm:"varchar(40)"`          //类型名称
	OpUser    string    `xorm:"index" json:"op_user"` //操作用户人
	CreatedAt time.Time `xorm:"created"`              //创建时间
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

//电影标签
type MovieTag struct {
	Id        int64
	Name      string    `xorm:"varchar(40)"`          //标签名称
	OpUser    string    `xorm:"index" json:"op_user"` //操作用户人
	CreatedAt time.Time `xorm:"created"`              //创建时间
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}
