/*
创建者：     Zero
创建时间：   2021/5/26
项目名称：   goWeb
*/
package model

import "gorm.io/gorm"

//用户信息
type UserInfo struct {
	gorm.Model
	ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`                 //用户唯一标识
	Name      string `gorm:"type:varchar(20);not null;index" json:"name"`        //用户名称
	Telephone string `gorm:"varchar(11);not null;unique;index" json:"telephone"` //用户手机号
	Password  string `gorm:"size:255;not null" json:"password"`                  //用户密码
	RoleId    int    `gorm:"not null;index;default 2" json:"role_id"`            //用户角色
	Profile   string `json:"profile"`                                            //用户头像
}

//用户角色
type UserRole struct {
	gorm.Model
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`          //角色唯一标识
	Name string `gorm:"type:varchar(20);not null;index" json:"name"` //角色名称
}

//用户角色对应用户权限
//用户登录操作日志
