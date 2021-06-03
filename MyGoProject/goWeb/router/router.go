/*
创建者：     Zero
创建时间：   2021/5/26
项目名称：   goWeb
*/
package router

import (
	"MyGoProject/controller"
	"MyGoProject/middleware"
	"github.com/gin-gonic/gin"
)

//最后可加一个https

func CollectRouter(r *gin.Engine) *gin.Engine {

	r.Use(middleware.CORSMiddleware(), middleware.RecoverMiddleware()) //使用中间件
	r.POST("/api/auth/login", controller.Login)
	r.POST("/api/auth/register", controller.Register)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	return r
}
