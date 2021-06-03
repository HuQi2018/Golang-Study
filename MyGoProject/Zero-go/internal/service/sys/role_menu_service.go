package sys

import (
	"Zero-go/internal/model/sys"
	"Zero-go/pkg/DB"
	"errors"
	"fmt"
)

var (
	SuperAdminRoutes = []string{"/sys/org_types", "/sys/roles", "/sys/menus", "/sys/role_menu", "/sys/users",
		"/sys/user_edit", "/sys/permission", "/sys/movies", "/sys/movie_tag", "/sys/movie_type"}
)

type RoleMenuService struct{}

const SELECT_ROLE_MENU = "select_role_menu.stpl"

type TempRoleMenu struct {
	Id        int64
	Checked   int
	RouteLink string
}

//此用户所拥有的所有页面路由查询
func (self *RoleMenuService) FindRouteLinksByRole(roleId, employeeId int64) (routeLinks []string, err error) {
	params := map[string]interface{}{
		"role_id": roleId,
		"checked": 1,
	}
	roleMenus := make([]*TempRoleMenu, 0)
	if err = DB.FindRowsBySqlTemplate(SELECT_ROLE_MENU, &params, &roleMenus); err != nil {
		return
	}
	//超管页面固定
	if roleId == 1 && employeeId == 1 {
		routeLinks = append(routeLinks, SuperAdminRoutes...)
		return routeLinks, err
	}
	for _, roleMenu := range roleMenus {
		routeLinks = append(routeLinks, roleMenu.RouteLink)
	}
	return routeLinks, err
}

// 根据角色查询所在组织的菜单权限，拥有权限的则勾选，不拥有的取消勾选
func (self *RoleMenuService) FindMenusByRole(roleId int64) ([]*sys.SysMenu, error) {
	role := sys.Role{}
	ok, err := DB.FindById(roleId, &role)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("非法的角色ID")
	}

	level1Menus := make([]*sys.SysMenu, 0)
	sql := "SELECT * FROM `sys_menu` where deleted_at is null and level=? and instr(org_type_ids, ?) order by `index` asc"
	err = DB.SQL(sql, "level1", role.OrgTypeId).FindRows(&level1Menus)
	if err != nil {
		return nil, err
	}

	level2Menus := make([]*sys.SysMenu, 0)
	err = DB.SQL(sql, "level2", role.OrgTypeId).FindRows(&level2Menus)
	if err != nil {
		return nil, err
	}

	roleMenus := make([]*sys.RoleMenu, 0)
	err = DB.Where("role_id=? and checked=1", roleId).FindRows(&roleMenus)
	if err != nil {
		return nil, err
	}

	for _, level2Menu := range level2Menus {
		for _, roleMenu := range roleMenus {
			if roleMenu.MenuId == level2Menu.Id {
				level2Menu.Checked = roleMenu.Checked
				break
			}
		}
	}

	for _, level1Menu := range level1Menus {
		level1Menu.Children = make([]*sys.SysMenu, 0)
		for _, level2Menu := range level2Menus {
			if level2Menu.ParentId == level1Menu.Id {
				level1Menu.Children = append(level1Menu.Children, level2Menu)
			}
		}
	}

	return level1Menus, err
}

//保存
func (self *RoleMenuService) Save(roleMenuForm *sys.RoleMenuReq) error {
	var (
		err error
		ok  bool
	)

	session := NewSession()
	if err = session.Begin(); err != nil {
		return err
	}
	defer session.Close()
	// 因为有角色取消权限和新增权限
	// 先把此角色的所有权限都设置为没有权限，然后重新分配权限
	sql := "update role_menu set checked = 0 where role_id = ?"
	if err = DB.ExecuteSqlTx(session, sql, roleMenuForm.RoleId); err != nil {
		return fmt.Errorf("修改role_menu checked状态报错: %s", err.Error())
	}

	roleMenus := make([]*sys.RoleMenu, 0)
	insertRoleMenus := make([]*sys.RoleMenu, 0)
	updateRoleMenus := make([]*sys.RoleMenu, 0)
	for _, menuId := range roleMenuForm.MenuIds {
		roleMenu := sys.RoleMenu{}
		if ok, err = DB.Where("menu_id = ? and role_id = ?", menuId, roleMenuForm.RoleId).FindOne(&roleMenu); err != nil {
			session.Rollback()
			return err
		}

		roleMenu.MenuId = menuId
		roleMenu.RoleId = roleMenuForm.RoleId
		roleMenu.Checked = 1
		if ok {
			updateRoleMenus = append(updateRoleMenus, &roleMenu)
		} else {
			insertRoleMenus = append(insertRoleMenus, &roleMenu)
		}
		roleMenus = append(roleMenus, &roleMenu)
	}

	// 新增角色权限
	if len(insertRoleMenus) > 0 {
		if _, err = DB.InsertOneTx(session, insertRoleMenus); err != nil {
			session.Rollback()
			return err
		}
	}

	// 修改角色权限
	if len(updateRoleMenus) > 0 {
		for _, m := range updateRoleMenus {
			if _, err = DB.UpdateByIdTx(session, m.Id, m); err != nil {
				session.Rollback()
				return err
			}
		}
	}

	err = session.Commit()
	return err
}
