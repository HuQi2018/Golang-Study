package sys

import "time"

type RolePermission struct {
	Id           int64
	RoleId       int64
	PermissionId int64
	Checked      int
	CreatedAt    time.Time `xorm:"created"`
	UpdatedAt    time.Time `xorm:"updated"`
	DeletedAt    time.Time `json:"-" xorm:"deleted"`
}

type RolePermissionReq struct {
	RoleId        int64   `json:"role_id"`
	PermissionIds []int64 `json:"permission_ids"`
}
