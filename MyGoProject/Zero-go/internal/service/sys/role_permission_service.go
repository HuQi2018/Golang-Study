package sys

import (
	"Zero-go/internal/model/sys"
	"Zero-go/pkg/DB"
	"errors"
	"fmt"
	"github.com/gogf/gf/g/os/glog"
)

type RolePermissionService struct{}

// 根据角色查询所在组织的权限，其中此角色拥有的权限则勾选，不拥有的取消勾选
func (self *RolePermissionService) FindPermissionByRole(roleId int64) ([]*sys.Permission, error) {
	role := sys.Role{}
	ok, err := DB.FindById(roleId, &role)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("非法的角色ID")
	}

	//找出符合条件的接口组
	groups := make([]*sys.Permission, 0)
	sql := "SELECT * FROM `permission` where deleted_at is null and group_id = ? and id not in (1,2,3) and instr(org_types, ?) order by `id` asc"
	if err = DB.SQL(sql, "0", role.OrgTypeId).FindRows(&groups); err != nil {
		return nil, err
	}
	//找到符合条件的接口
	perms := make([]*sys.Permission, 0)
	sql = "SELECT * FROM `permission` where deleted_at is null and group_id != ? and instr(org_types, ?) order by `id` asc"
	if err = DB.SQL(sql, "0", role.OrgTypeId).FindRows(&perms); err != nil {
		return nil, err
	}

	//查询此角色当前被赋予的权限
	rolePermissions := make([]*sys.RolePermission, 0)
	if err = DB.Where("role_id = ? and checked = 1", roleId).FindRows(&rolePermissions); err != nil {
		return nil, err
	}
	//复制 role_permission表中的勾选情况
	for _, perm := range perms {
		for _, rolePermission := range rolePermissions {
			if rolePermission.PermissionId == perm.Id {
				perm.Checked = rolePermission.Checked
				break
			}
		}
	}

	//得到此角色所对应的所有权限
	for _, group := range groups {
		group.Children = make([]*sys.Permission, 0)
		for _, perm := range perms {
			if perm.GroupId == group.Id {
				group.Children = append(group.Children, perm)
			}
		}
	}

	return groups, err
}

//保存
func (self *RolePermissionService) Save(rolePermissionReq *sys.RolePermissionReq) error {
	var (
		err error
		ok  bool
	)

	session := NewSession()
	if err = session.Begin(); err != nil {
		return err
	}
	defer session.Close()

	//找到此角色
	role := sys.Role{}
	ok, err = DB.FindById(rolePermissionReq.RoleId, &role)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("非法的角色ID")
	}

	// 先把该角色的所有权限都设置为0，然后重新分配权限
	sql := "update role_permission set checked = 0 where role_id = ?"
	if err = DB.ExecuteSqlTx(session, sql, rolePermissionReq.RoleId); err != nil {
		return fmt.Errorf("修改role_permission checked状态报错: %s", err.Error())
	}

	//查询此角色所属层级所有权限
	perms := make([]*sys.Permission, 0)
	sql = "SELECT * FROM `permission` where deleted_at is null and group_id != ? and instr(org_types, ?) order by `id` asc"
	if err = DB.SQL(sql, "0", role.OrgTypeId).FindRows(&perms); err != nil {
		return err
	}
	//更新所有权限的状态
	for _, perm := range perms {
		for _, permissionId := range rolePermissionReq.PermissionIds {
			if perm.Id != permissionId {
				continue
			}
			perm.Checked = 1
			break
		}
	}

	//role_permission表中已存在的权限
	updatePerms := make([]*sys.RolePermission, 0)
	sql = "SELECT * FROM `role_permission` where deleted_at is null and role_id = ?"
	if err = DB.SQL(sql, role.Id).FindRows(&updatePerms); err != nil {
		return err
	}

	//对比来找到需要被新增的权限
	var newRolePerms []*sys.RolePermission
	var flag = false
	for _, perm := range perms {
		for _, updatePerm := range updatePerms {
			if perm.Id != updatePerm.PermissionId {
				continue
			}
			//如果数据库中存在此记录，则更新状态
			updatePerm.Checked = perm.Checked
			if _, err = DB.UpdateByIdTx(session, updatePerm.Id, updatePerm); err != nil {
				glog.Error(err)
				return err
			}
			flag = true
			break
		}
		if !flag { //如果没有此perm记录，则新增
			rolePerm := sys.RolePermission{}
			rolePerm.RoleId = role.Id
			rolePerm.PermissionId = perm.Id
			rolePerm.Checked = perm.Checked
			newRolePerms = append(newRolePerms, &rolePerm)
		}
		flag = false
	}

	//把所有新的权限新增到role_permission表中
	if _, err = DB.InsertBatchTx(session, newRolePerms); err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	return err
}
