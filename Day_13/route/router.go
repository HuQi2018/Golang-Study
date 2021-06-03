/*
创建者：     Zero
创建时间：   2021/5/25
项目名称：   golang-study
*/
package route

import (
	"github.com/gin-gonic/gin"
	"golang-study/huqi/Day_13/controller"
	"golang-study/huqi/Day_13/middleware"
)

func CollectRouter(r *gin.Engine) *gin.Engine {

	r.Use(middleware.CORSMiddleware(), middleware.RecoverMiddleware()) //使用中间件
	r.POST("/api/auth/login", controller.Login)
	r.POST("/api/auth/register", controller.Register)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	return r
}
