package v1

import (
	model "Zero-go/internal/model/sys"
	service "Zero-go/internal/service/sys"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/g/os/glog"
)

type PermissionController struct {
	Router gin.IRouter
}

var permissionService service.PermissionService

func (self *PermissionController) Setup() {
	self.Router.GET("/permission/query", self.Find)
	self.Router.POST("/permission/save", self.Save)
	self.Router.GET("/permission/del", self.Delete)
}

func (*PermissionController) Find(c *gin.Context) {
	page, limit := GetPageParams(c)
	keyword := c.Query("key")
	perms, err := permissionService.Find(page, limit, keyword)
	if err != nil {
		glog.Error(err)
		ResponseError(c, err)
		return
	}

	ResponseOK(c, perms)
}

func (*PermissionController) Save(c *gin.Context) {
	req := model.PermissionReq{}
	if err := BindJSON(c, &req); err != nil {
		glog.Error(err)
		ResponseErrorf(c, "请求的数据不合法: %s", err)
		return
	}
	if err := permissionService.Save(&req); err != nil {
		glog.Error(err)
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}

func (*PermissionController) Delete(c *gin.Context) {
	id := GetId(c)
	if err := permissionService.Delete(id); err != nil {
		glog.Error(err)
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}
