package permission_check

import (
	"Zero-go/internal/common"
	"Zero-go/internal/common/middleware/jwt"
	"Zero-go/internal/common/permission_check"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/g/os/glog"
)

func jwtTokenAbort(c *gin.Context, code int, msg string) {
	common.ResponseWith(c, code, msg)
	c.Abort()
}

//权限检查中间件
func AuthCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		//根据上下文获取载荷claims 从claims获得role
		claims := c.MustGet("userClaims").(*jwt.UserClaims)
		//检查权限
		permissionCheck := permission_check.PermissionCheckModel{
			RoleId:     claims.RoleId,
			EmployeeId: claims.UserId,
			Url:        c.Request.URL.Path,
		}
		flag, apiName, err := permissionCheck.CheckPermission()
		if err != nil {
			glog.Error(err)
			jwtTokenAbort(c, 500, err.Error())
			return
		}
		if !flag {
			errMsg := fmt.Sprintf("很抱歉您没有此权限 [%s]", apiName)
			jwtTokenAbort(c, 500, errMsg)
			return
		}
		c.Next()
	}
}
