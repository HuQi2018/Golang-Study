package sys

import "time"

type Permission struct {
	Id        int64
	GroupId   int64         `json:"group_id"`
	Name      string        `xorm:"varchar(32) comment('组或接口的名称')" json:"name"`
	Url       string        `xorm:"varchar(64) comment('接口地址')" json:"url"`
	OrgTypes  string        `xorm:"varchar(256) comment('接口地址')" json:"org_types"`
	BuildIn   int           `xorm:"int(2) comment('是否为超管固有权限')" json:"build_in"`
	ChildNum  int           `xorm:"int(16) comment('组内接口数量')" json:"child_num"`
	Children  []*Permission `xorm:"-"`
	Checked   int           `xorm:"-"`
	CreatedAt time.Time     `xorm:"created"`
	UpdatedAt time.Time     `xorm:"updated"`
	DeletedAt time.Time     `json:"-" xorm:"deleted"`
}

type PermissionReq struct {
	Id       int64
	GroupId  int64  `json:"group_id"`
	Name     string `xorm:"varchar(32) comment('组或接口的名称')" json:"name"`
	Url      string `xorm:"varchar(64) comment('接口地址')" json:"url"`
	OrgTypes []int  `xorm:"varchar(256) comment('接口地址')" json:"org_types"`
}
