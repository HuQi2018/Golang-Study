/**
 * Created by Wangwei on 2019-06-04 17:21.
 */

package sys

import (
	"Zero-go/internal/conf"
	"Zero-go/internal/model/dto"
	"Zero-go/internal/model/sys"
	"Zero-go/pkg/DB"
	"Zero-go/pkg/util/gosecurity"
	"errors"
	"github.com/xormplus/xorm"
)

var (
	SELECT_USER_TOTAL = "select_user_total.stpl"
	SELECT_USER       = "select_user.stpl"
)

type UserService struct{}

func (self *UserService) Get(id int64) (*sys.User, error) {
	user := sys.User{}
	has, err := DB.FindById(id, &user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("账户不存在")
	}

	return &user, err
}

//分页查询
func (self *UserService) FindByPage(page, limit int, keyword string, orgTypeId, orgId int64, isAdmin int) (pages *dto.Pages, err error) {
	params := map[string]interface{}{
		"offset":      GetOffset(page, limit),
		"limit":       limit,
		"keyword":     LikeStr(keyword),
		"org_type_id": orgTypeId,
		"org_id":      orgId,
		"is_admin":    isAdmin,
	}

	var total int64
	user := make([]*sys.User, 0)
	err = DB.PageBySqlTemplateClient(SELECT_USER, &params, &user, SELECT_USER_TOTAL, &total)

	pages = &dto.Pages{Total: total, Data: user}
	return pages, err
}

func (UserService) GetByAccount(account, password string) (user *sys.User, err error) {
	var ok bool
	user = &sys.User{}
	ok, err = DB.Where("`account`=? and `password`=?", account, password).Omit("password").FindOne(user)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("用户名/密码错误")
	}
	return
}

//通过id查询
func (self *UserService) FindById(id int64) (*sys.User, error) {
	user := new(sys.User)
	ok, err := DB.FindById(id, user)
	if !ok {
		return nil, errors.New("员工ID:" + string(user.Id) + "不存在")
	}

	return user, err
}

func (self *UserService) Save(user *sys.User) (err error) {
	session := NewSession()
	defer session.Close()

	return self.SaveTx(session, user)
}

func (UserService) SaveTx(session *xorm.Session, user *sys.User) (err error) {
	count, err := DB.Where("account = ?", user.Account).Count(&sys.User{})
	if err != nil {
		return err
	}

	if user.Id == 0 && count > 0 {
		return errors.New("账号已存在，请检查！")
	}

	if user.Id == 0 {
		user.Password = gosecurity.MD5Password(user.Password)
		_, err = DB.InsertOneTx(session, user)
	} else {
		_, err = DB.UpdateByIdWithOmitTx(session, user.Id, user, "account", "password")
	}

	return
}

func (self *UserService) ResetPassword(account string) (err error) {

	user := new(sys.User)
	ok, err := DB.Where("`account`=? ", account).FindOne(user)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("无此用户")
	}

	user.Password = gosecurity.MD5Password(conf.AppConf.DefaultPassword)
	_, err = DB.UpdateById(user.Id, user)

	return err
}

func (self *UserService) UpdatePassword(req *sys.PasswordReq) (err error) {

	user := new(sys.User)
	ok, err := DB.Where("`account`=? and `password`=?", req.Account, gosecurity.MD5Password(req.OldPassword)).FindOne(user)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("密码错误")
	}

	user.Password = gosecurity.MD5Password(req.Password)
	_, err = DB.UpdateById(user.Id, user)

	return err
}
