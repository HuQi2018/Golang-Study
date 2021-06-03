package permission_check

import (
	"Zero-go/internal/model/sys"
	"Zero-go/pkg/DB"
)

type PermissionCheckModel struct {
	RoleId     int64
	EmployeeId int64
	Url        string
}

type TempPermission struct {
	Id      int64
	Checked int
	BuildIn int
	Name    string
}

var (
	SELECT_PERMISSION_CHECK = "select_permission_check.stpl"
)

func (self *PermissionCheckModel) CheckPermission() (flag bool, apiName string, err error) {

	//超管直接不检查
	if self.EmployeeId == 1 {
		return true, "", err
	}

	//判断是否为超管， 请求的接口是否为超管固有权限。若超管请求非固有接口，则拒绝请求
	//此处设计为超管只能使用超管页面的接口，即build_in为 1 的“固有接口”
	//if self.EmployeeId == 1 {
	//	perm := sys.Permission{}
	//	has, err := DB.Where("url=?", self.Url).FindOne(&perm)
	//	if err != nil {
	//		return false, "", err
	//	}
	//	if !has {
	//		return false, "", fmt.Errorf("未找到此权限 [%s]", self.Url)
	//	}
	//	if perm.BuildIn == 1 {
	//		return true, "", err
	//	}
	//	return false, perm.Name, err
	//}

	params := map[string]interface{}{
		"role_id": self.RoleId,
		"url":     self.Url,
	}

	tempPermission := TempPermission{}
	has, err := DB.FindOneBySqlTemplate(SELECT_PERMISSION_CHECK, &params, &tempPermission)
	if err != nil {
		return false, "", err
	}

	//正常判断权限
	if has && tempPermission.Checked == 1 {
		return true, "", err
	}
	//如果role_permission表中根本没有此数据，说明在权限管理的时候层级并没有选中他所在的层级
	//此时则去查询一次此接口的名称，便于管理员去添加对应的层级
	if len(tempPermission.Name) == 0 {
		perm := sys.Permission{}
		if has, err = DB.Where("url=?", self.Url).FindOne(&perm); err != nil {
			return false, "", err
		}
		if !has {
			return false, "", err
		}
		tempPermission.Name = perm.Name
	}
	return false, tempPermission.Name, err
}
