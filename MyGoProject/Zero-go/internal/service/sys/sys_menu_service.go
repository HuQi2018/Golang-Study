package sys

import (
	"Zero-go/internal/model/sys"
	"Zero-go/pkg/DB"
	"errors"

	"github.com/xormplus/xorm"
)

type SysMenuService struct{}

func (self *SysMenuService) GetPermsByApiUrl(apiUrl string) (*sys.SysMenu, error) {
	sysMenu := new(sys.SysMenu)
	has, err := DB.Where("api_url=?", apiUrl).FindOne(sysMenu)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("未找到此记录")
	}
	return sysMenu, err
}

func (self *SysMenuService) FindAll() ([]*sys.SysMenu, error) {
	level1Menus := make([]*sys.SysMenu, 0)
	err := DB.Where("level=?", "level1").Asc("index").FindRows(&level1Menus)
	if err != nil {
		return nil, err
	}

	level2Menus := make([]*sys.SysMenu, 0)
	err = DB.Where("level=?", "level2").Asc("index").FindRows(&level2Menus)
	if err != nil {
		return nil, err
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
func (self *SysMenuService) Save(sysMenu *sys.SysMenu) (err error) {
	if len(sysMenu.OrgTypeIds) == 0 {
		sysMenu.OrgTypeIds = []int64{0}
	}

	if sysMenu.Id == 0 {
		_, err = DB.InsertOne(sysMenu)
		return err
	}

	session := NewSession()
	defer session.Close()
	if err = session.Begin(); err != nil {
		return err
	}

	if _, err = DB.UpdateByIdTx(session, sysMenu.Id, sysMenu); err != nil {
		session.Rollback()
		return err
	}

	if err = self.UpdateChilds(session, sysMenu); err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	return
}

func (self *SysMenuService) UpdateChilds(session *xorm.Session, sysMenu *sys.SysMenu) (err error) {
	if sysMenu.Level != "level2" {
		childMenus := make([]*sys.SysMenu, 0)
		if err = DB.WhereTx(session, "parent_id=?", sysMenu.Id).FindRows(&childMenus); err != nil {
			return err
		}

		if len(childMenus) > 0 {
			for _, child := range childMenus {
				child.OrgTypeIds = sysMenu.OrgTypeIds
				if _, err = DB.UpdateByIdTx(session, child.Id, child); err != nil {
					return err
				}
			}
		}
	}

	return
}

// 删除菜单
func (self *RoleMenuService) Delete(menuId int64) error {
	session := NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	// 递归调用，从当前菜单的子菜单最后一级开始删除，然后依次删除上级菜单
	if err := self.deleteChilds(session, menuId); err != nil {
		session.Rollback()
		return err
	}

	if err := session.Commit(); err != nil {
		return err
	}

	return nil
}

// 递归调用，从当前菜单的子菜单最后一级开始删除，然后依次删除上级菜单
func (self *RoleMenuService) deleteChilds(session *xorm.Session, parentId int64) error {
	var (
		ok         bool
		err        error
		parentMenu sys.SysMenu
	)

	ok, err = DB.FindById(parentId, &parentMenu)
	if !ok {
		return errors.New("菜单ID不存在")
	}
	if err != nil {
		return err
	}

	if parentMenu.Level == "level1" {
		// 当前菜单为level1级，先查找下面所有level2级菜单
		level2Menus := make([]*sys.SysMenu, 0)
		if err = DB.Where("parent_id = ?", parentId).FindRows(&level2Menus); err != nil {
			return err
		}

		for _, level2Menu := range level2Menus {
			// level2级菜单，直接删除
			if _, err = session.ID(level2Menu.Id).Delete(level2Menu); err != nil {
				goto ERR
			}
		}
	}

	// 最后删除当前菜单，不管当前菜单是level1, level2 他们的子菜单都已经被上面步骤删除了
	if _, err = session.ID(parentId).Delete(&parentMenu); err != nil {
		return err
	}

ERR:
	return err
}
