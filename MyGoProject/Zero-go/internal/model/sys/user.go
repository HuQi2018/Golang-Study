package sys

import (
	"time"
)

type User struct {
	Id        int64     `json:"id"`
	Code      string    `xorm:"varchar(50) comment('用户编码')" json:"code"`
	Type      string    `xorm:"varchar(50) comment('用户类型 admin管理员 user普通用户')" json:"type"`
	Name      string    `xorm:"varchar(80) notnull" json:"name" binding:"required"`
	Account   string    `xorm:"varchar(40) notnull" json:"account" binding:"required"`
	Password  string    `xorm:"varchar(80) notnull" json:"password,omitempty" binding:"required"`
	Phone     string    `xorm:"varchar(11)" json:"phone" binding:"phone"`
	RoleId    int64     `json:"role_id"`
	RoleName  string    `xorm:"varchar(80)" json:"role_name"`
	OrgTypeId int64     `xorm:"comment('机构类型ID')"`
	OrgId     int64     `xorm:"comment('机构ID')"`
	OrgName   string    `xorm:"varchar(40) comment('机构名称')"`
	Profile   string    `xorm:"comment('用户头像')" json:"profile"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
	DeletedAt time.Time `xorm:"deleted" json:"deleted_at"`
}

type LoginInfoReq struct {
	//Code     string `json:"code" binding:"required"`
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	//Type     string
}

type PasswordReq struct {
	Account     string
	OldPassword string
	Password    string
	Password2   string
}

type UserReq struct {
	Account string
}
