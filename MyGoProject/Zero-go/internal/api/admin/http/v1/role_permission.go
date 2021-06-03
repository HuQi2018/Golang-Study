package v1

import (
	"Zero-go/internal/model/sys"
	sys_service "Zero-go/internal/service/sys"
	"github.com/gogf/gf/g/os/glog"

	"github.com/gin-gonic/gin"
)

var (
	rolePermissionService = &sys_service.RolePermissionService{}
)

type RolePermissionController struct {
	Router gin.IRouter
}

func (self *RolePermissionController) Setup() {
	self.Router.GET("/role_permission/query", self.FindAll)
	self.Router.POST("/role_permission/save", self.Save)
}

func (self *RolePermissionController) FindAll(c *gin.Context) {
	roleId := GetInt64("role_id", c)
	rolePermissions, err := rolePermissionService.FindPermissionByRole(roleId)
	if err != nil {
		glog.Error(err)
		ResponseError(c, err)
		return
	}

	ResponseOK(c, rolePermissions)
}

func (self *RolePermissionController) Save(c *gin.Context) {
	rolePermissionReq := sys.RolePermissionReq{}
	if err := BindJSON(c, &rolePermissionReq); err != nil {
		glog.Error(err)
		ResponseErrorf(c, "请求的数据不合法: %s", err)
		return
	}
	if err := rolePermissionService.Save(&rolePermissionReq); err != nil {
		glog.Error(err)
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}
