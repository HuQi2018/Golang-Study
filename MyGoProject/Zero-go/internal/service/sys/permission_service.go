package sys

import (
	"Zero-go/internal/model/dto"
	"Zero-go/internal/model/sys"
	"Zero-go/pkg/DB"
	"Zero-go/pkg/util/goconv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"strings"
)

type PermissionService struct{}

var (
	SELECT_PERMISSION_TOTAL = "select_permission_total.stpl"
	SELECT_PERMISSION       = "select_permission.stpl"
)

func (*PermissionService) Find(page, limit int, keyword string) (pages *dto.Pages, err error) {
	results := make([]*sys.Permission, 0)

	// 找到所有组
	groups := []*sys.Permission{}
	sql := "group_id = 0"
	if len(keyword) > 0 {
		sql += fmt.Sprintf(" and name like '%%%s%%'", keyword)
	}
	if err = DB.Where(sql).FindRows(&groups); err != nil {
		return
	}

	// 前端表格只有在收到若干个完整的组的时候才能正确显示出来
	// 根据limit来限制返回的组数（接口数量） limit >= 组数量 + 接口数量
	num := 0
	tempNum := 0
	offset := GetOffset(page, limit)
	maxGroupId := int64(1)
	minGroupId := int64(0)
	for _, group := range groups {
		// num + 接口组记录（1） + 接口组内接口数量（group.ChildNum）
		if limit >= num+group.ChildNum+1 {
			tempNum += group.ChildNum + 1
			// 判断是否大于偏移量
			if tempNum <= offset {
				minGroupId = group.Id
				continue
			}
			num += group.ChildNum + 1
			if group.Id > maxGroupId {
				maxGroupId = group.Id
			}
			continue
		}
		break
	}

	params := map[string]interface{}{
		"keyword":      LikeStr(keyword),
		"max_group_id": maxGroupId,
		"min_group_id": minGroupId,
	}

	var total int64
	perms := make([]*sys.Permission, 0)
	if err = DB.PageBySqlTemplateClient(SELECT_PERMISSION, &params, &perms, SELECT_PERMISSION_TOTAL, &total); err != nil {
		return
	}

	// 先把接口组都放进去
	for _, group := range groups {
		if group.Id <= maxGroupId && group.Id > minGroupId {
			results = append(results, group)
		}
	}
	// 再把接口放进对应的接口组里
	for _, group := range results {
		for _, perm := range perms {
			if perm.GroupId == group.Id {
				group.Children = append(group.Children, perm)
			}
		}

	}

	pages = &dto.Pages{total, &results}
	return
}

func (*PermissionService) Save(permReq *sys.PermissionReq) (err error) {
	perm := sys.Permission{}
	if err = copier.Copy(&perm, permReq); err != nil {
		return
	}
	orgTypesBytes, err := json.Marshal(permReq.OrgTypes)
	if err != nil {
		return err
	}
	orgTypes := string(orgTypesBytes)[1 : len(string(orgTypesBytes))-1]
	perm.OrgTypes = orgTypes

	session := NewSession()
	if err := session.Begin(); err != nil {
		return err
	}
	defer session.Close()

	// 新增
	if perm.Id == 0 {
		perm.ChildNum = 0
		if _, err = DB.InsertOneTx(session, &perm); err != nil {
			return err
		}
		// 如果是接口，则将对应的接口组的childNum +1
		if perm.GroupId != 0 {
			group := sys.Permission{}
			has, err := DB.FindById(perm.GroupId, &group)
			if err != nil {
				return err
			}
			if !has {
				return errors.New("未查询到所属接口组")
			}
			group.ChildNum++
			if _, err = DB.UpdateByIdTx(session, group.Id, &group); err != nil {
				return err
			}
		}
		err = session.Commit()
		return err
	}

	// 更新
	if _, err = DB.UpdateByIdTx(session, perm.Id, &perm); err != nil {
		return err
	}

	// 判断所属机构是否发生变化，而且变化如 *原本属于某层级的权限被取消了*
	// 那么需要去把数据库中这个层级的所有角色关于这个权限的记录给删除掉
	// 以避免用户在修改了菜单之后却不更新权限，导致此用户将继续拥有此权限的情况
	oldPerm := sys.Permission{}
	has, err := DB.FindById(perm.Id, &oldPerm)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("未找到此记录")
	}
	if oldPerm.OrgTypes == perm.OrgTypes { // 如果没变化，则结束
		err = session.Commit()
		return err
	}

	// 层级发生了变化，判断变化状态
	oldOrgTypes := strings.Split(oldPerm.OrgTypes, ",")
	var deletedOrgTypes = []int{}
	var flag = 0
	for _, oldOrgType := range oldOrgTypes {
		tOldOrgType := goconv.ToInt(oldOrgType)
		for _, orgType := range permReq.OrgTypes {
			if orgType == tOldOrgType {
				flag = tOldOrgType
				break
			}
		}
		if flag == 0 {
			deletedOrgTypes = append(deletedOrgTypes, tOldOrgType)
		}
		flag = 0
	}

	// 如果被删除的层级数为0, 则结束
	if len(deletedOrgTypes) == 0 {
		err = session.Commit()
		return err
	}

	// 层级数不为0， 则找到这些层级下的所有role_id
	roleIds := make([]int64, 0)
	for _, orgTypeId := range deletedOrgTypes {
		tRoles := make([]*sys.Role, 0)
		if err = DB.Where("org_type_id = ?", orgTypeId).FindRows(&tRoles); err != nil {
			return err
		}
		for _, role := range tRoles {
			roleIds = append(roleIds, role.Id)
		}
	}

	// 找到role_permission表中，permission_id为刚刚更新的id，且role_id属于被删除的层级的记录
	rolePermissions := make([]*sys.RolePermission, 0)
	for _, roleId := range roleIds {
		rolePermission := sys.RolePermission{}
		has, err := DB.Where("permission_id = ? and role_id = ?", perm.Id, roleId).FindOne(&rolePermission)
		if err != nil {
			return err
		}
		if has {
			rolePermissions = append(rolePermissions, &rolePermission)
		}
	}

	// 删除
	for _, rolePermission := range rolePermissions {
		if _, err = DB.DeleteByIdTx(rolePermission.Id, rolePermission, session); err != nil {
			return err
		}
	}
	err = session.Commit()
	return err
}

func (*PermissionService) Delete(id int64) (err error) {
	perm := sys.Permission{}
	has, err := DB.FindById(id, &perm)
	if err != nil {
		return
	}
	if !has {
		return errors.New("未找到此记录")
	}

	if perm.GroupId != 0 {
		_, err = DB.DeleteById(id, &perm)
		return
	}
	// 接口组级联删除
	session := NewSession()
	if err = session.Begin(); err != nil {
		return
	}
	defer session.Close()
	perms := make([]*sys.Permission, 0)
	if err = DB.Where("group_id = ?", perm.Id).FindRows(&perms); err != nil {
		return
	}
	// 删除组内所有接口
	for _, p := range perms {
		if _, err = DB.DeleteById(p.Id, p); err != nil {
			session.Rollback()
			return
		}
	}
	// 删除接口组
	if _, err = DB.DeleteById(id, &perm); err != nil {
		session.Rollback()
		return
	}
	err = session.Commit()
	return
}
