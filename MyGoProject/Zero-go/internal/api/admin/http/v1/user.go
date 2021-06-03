/**
 * Created by Wangwei on 2019-06-03 20:09.
 */

package v1

import (
	"Zero-go/internal/common/enum/employee_enum"
	"Zero-go/internal/model/sys"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Router gin.IRouter
}

func (self *UserController) Setup() {
	self.Router.GET("/user/get", self.Get)
	self.Router.GET("/user/query", self.Find)
	self.Router.POST("/user/save", self.Save)
	self.Router.POST("/user/reset_password", self.ResetPassword)
	self.Router.POST("/user/update_password", self.UpdatePassword)
}

func (self *UserController) Find(c *gin.Context) {
	claims := GetUserClaims(c)
	page, limit := GetPageParams(c)
	keyword := c.Query("keyword")

	var isAdmin int
	if claims.RoleId == 1 {
		isAdmin = 1
	}

	user, err := userService.FindByPage(page, limit, keyword, claims.OrgTypeId, claims.OrgId, isAdmin)
	if err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, user)
}

func (self *UserController) Get(c *gin.Context) {
	userId := GetId(c)
	user, err := userService.FindById(userId)
	if err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, user)
}

func (self *UserController) Save(c *gin.Context) {
	user := new(sys.User)
	if err := BindJSON(c, user); err != nil {
		ResponseError(c, err)
		return
	}

	claims := GetUserClaims(c)
	user.OrgTypeId = claims.OrgTypeId
	user.OrgId = claims.OrgId
	user.OrgName = claims.OrgName
	switch user.RoleId {
	case 1:
		user.Type = employee_enum.SuperAdmin
	case 2:
		user.Type = employee_enum.ADMIN
	case 3:
		user.Type = employee_enum.USER
	}
	if err := userService.Save(user); err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}

func (self *UserController) ResetPassword(c *gin.Context) {

	userReq := new(sys.UserReq)
	if err := BindJSON(c, userReq); err != nil {
		ResponseError(c, err)
		return
	}

	if err := userService.ResetPassword(userReq.Account); err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}

func (self *UserController) UpdatePassword(c *gin.Context) {

	passwordReq := new(sys.PasswordReq)
	if err := BindJSON(c, passwordReq); err != nil {
		ResponseError(c, err)
		return
	}

	if err := userService.UpdatePassword(passwordReq); err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}
