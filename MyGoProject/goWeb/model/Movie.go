/*
创建者：     Zero
创建时间：   2021/5/26
项目名称：   goWeb
*/
package model

import "gorm.io/gorm"

//电影信息
type MovieInfo struct {
	gorm.Model
	ID            int    `gorm:"primaryKey;autoIncrement" json:"id"`                 //唯一表标识
	MovieId       int    `gorm:"uniqueIndex" json:"movie_id"`                        //豆瓣电影唯一标识
	Title         string `gorm:"type:varchar(1000);not null;index" json:"name"`      //电影名称
	Genres        string `gorm:"type:varchar(20);not null;" json:"type"`             //电影类型
	Rating        string `gorm:"type:longtext;not null;" json:"rating"`              //电影评分
	Durations     string `gorm:"type:varchar(20);not null;" json:"durations"`        //时长
	Year          int    `gorm:"type:varchar(20);not null;" json:"year"`             //上映年份
	Pubdates      string `gorm:"type:varchar(1000);not null;" json:"pubdates"`       //上映日期数据
	Pubdate       string `gorm:"type:varchar(1000);not null;" json:"pubdate"`        //上映日期
	Aka           string `grom:"type:varchar(1000);not null;" json:"aka"`            //又名
	Tags          string `gorm:"type:varchar(1000);not null;" json:"tags"`           //标签
	OriginalTitle string `gorm:"type:varchar(1000);not null;" json:"original_title"` //原始标题
	Language      string `gorm:"type:varchar(20);not null;" json:"language"`         //电影语言
	Country       string `gorm:"type:varchar(20);not null;" json:"country"`          //制片国家
	Actors        string `gorm:"type:varchar(20);not null;" json:"actor"`            //演员
	Writers       string `gorm:"type:varchar(20);not null;" json:"writers"`          //作者
	Directors     string `gorm:"type:varchar(20);not null;" json:"directors"`        //导演
	Summary       string `gorm:"type:varchar(20);not null;" json:"summary"`          //简介
	Photos        string `gorm:"type:varchar(20);not null;" json:"photos"`           //图像数据
	Images        string `gorm:"type:varchar(20);not null;" json:"images"`           //海报数据
	Videos        string `gorm:"type:varchar(20);not null;" json:"videos"`           //视频数据
}
