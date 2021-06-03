/*
创建者：     Zero
创建时间：   2021/5/25
项目名称：   golang-study
*/
package model

import "gorm.io/gorm"

type UserBase struct {
	gorm.Model
	ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"type:varchar(20);not null" json:"name"`
	Telephone string `gorm:"varchar(11);not null;unique" json:"telephone"`
	Password  string `gorm:"size:255;not null" json:"password"`
}
