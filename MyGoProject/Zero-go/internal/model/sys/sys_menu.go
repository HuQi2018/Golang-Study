package sys

import "time"

type SysMenu struct {
	Id         int64
	ParentId   int64   `xorm:"comment('上级节点id')"`
	Name       string  `xorm:"varchar(20) comment('节点名称')"`
	RouteLink  string  `xorm:"varchar(200) comment('前端页面路由')"`
	OrgTypeIds []int64 `xorm:"comment('节点所属组织类型，可多选')"`
	Level      string  `xorm:"varchar(20) comment('节点级别')"`
	Index      int
	Children   []*SysMenu `xorm:"-"` //不创建表字段，仅填充使用
	RoleId     int64      `xorm:"-"` //不创建表字段，仅填充使用
	Checked    int        `xorm:"-"` //不创建表字段，仅填充使用
	NodeType   string     `json:"node_type"`
	ApiUrls    string     `xorm:"TEXT" json:"api_urls"`
	CreatedAt  time.Time  `xorm:"created"`
	UpdatedAt  time.Time  `xorm:"updated"`
	DeletedAt  time.Time  `json:"-" xorm:"deleted"`
}

type ApiUrl struct {
	ApiUrl string
}
