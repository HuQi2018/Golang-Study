package model

import (
	"Zero-go/internal/common/enum/employee_enum"
	"Zero-go/internal/common/enum/role_type_enum"
	"Zero-go/internal/conf"
	"Zero-go/internal/model/sys"
	"Zero-go/pkg/DB"
	"Zero-go/pkg/util/gosecurity"
	"fmt"
	"github.com/gogf/gf/g/os/glog"
)

/*
	数据库表绑定，同时初始化，若无表则初始化表
	同时初始化超级管理员账号信息
*/

func Init() {
	//需要同步的表结构
	if err := DB.Engine.Sync2(
		new(sys.User),
		new(sys.Role),
		new(sys.OrgType),
		new(sys.SysMenu),
		new(sys.RoleMenu),
		new(sys.Permission),
		new(sys.RolePermission),
		new(sys.MovieInfo),
		new(sys.MovieType),
		new(sys.MovieTag),
	); err != nil {
		panic(err)
	}

	// 初始化数据
	if err := initData(); err != nil {
		panic(err)
	}
}

func initData() (err error) {
	if err = initRole(); err != nil {
		return err
	}

	if err = initAccount(); err != nil {
		return err
	}

	return
}

func initRole() (err error) {
	count, err := DB.Where("id = 1").Count(&sys.Role{})
	if err != nil {
		return fmt.Errorf("init superadmin role err: %v\n", err)
	}

	if count > 0 {
		return nil
	}

	role := &sys.Role{Id: 1, Code: "1001", IsAdmin: 1, Name: "超级管理员", Buildin: 1}
	_, err = DB.InsertOne(role)
	return err
}

// 初始化超级管理员账号
func initAccount() (err error) {
	count, err := DB.Where("account=?", conf.AppConf.SuperAdminAccount).Count(&sys.User{})
	if err != nil {
		glog.Fatalf("init superadmin account err: %v\n", err)
		panic(err)
	}
	if count > 0 {
		return
	}

	session := DB.Engine.NewSession()
	defer session.Close()

	if err = session.Begin(); err != nil {
		return fmt.Errorf("session begin err: %s", err)
	}

	password := gosecurity.MD5Password(conf.AppConf.SuperAdminPassord)
	employee := &sys.User{
		Id:       1,
		Name:     "超级管理员",
		Account:  conf.AppConf.SuperAdminAccount,
		Password: password,
		Phone:    conf.AppConf.SuperAdminPhone,
		RoleId:   role_type_enum.SuperAdmin,
		Code:     "1000",
		RoleName: "超级管理员",
		Type:     employee_enum.SuperAdmin,
	}
	if _, err = DB.InsertOneTx(session, employee); err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	return err
}
