/**
 * Created by Wangwei on 2019-06-05 11:26.
 */

package v1

import (
	"Zero-go/internal/common/middleware/jwt"
	"Zero-go/internal/common/middleware/permission_check"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	router := r.Group("/v1/admin_api")
	SetupNoneAuthorized(router)

	router.Use(jwt.UserJWTAuth())
	router.Use(permission_check.AuthCheckRole())
	SetupAuthorized(router)
}

// 不需要token认证的接口
func SetupNoneAuthorized(router gin.IRouter) {
	authController := AuthController{}
	router.POST("/login", authController.Login)
}

// 需要token认证的接口
func SetupAuthorized(router gin.IRouter) {
	orgTypeController := OrgTypeController{router}
	orgTypeController.Setup()

	roleController := RoleController{router}
	roleController.Setup()

	roleMenuController := RoleMenuController{router}
	roleMenuController.Setup()

	sysMenuController := SysMenuController{router}
	sysMenuController.Setup()

	userController := UserController{router}
	userController.Setup()

	permissionController := PermissionController{router}
	permissionController.Setup()

	rolePermissionController := RolePermissionController{router}
	rolePermissionController.Setup()

	movieController := MovieController{router}
	movieController.Setup()

	movieTagController := MovieTagController{router}
	movieTagController.Setup()

	movieTypeController := MovieTypeController{router}
	movieTypeController.Setup()

}
