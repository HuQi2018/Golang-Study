/*
创建者：     Zero
创建时间：   2021/5/28
项目名称：   Zero-go
*/
package v1

import (
	"Zero-go/internal/model/sys"
	sys_service "Zero-go/internal/service/sys"
	"github.com/gin-gonic/gin"
	"strconv"
)

var (
	movieTypeService = &sys_service.MovieTypeService{}
)

func (self *MovieTypeController) Setup() {
	self.Router.GET("/movie_type/get", self.Get)
	self.Router.GET("/movie_type/query", self.Find)
	self.Router.POST("/movie_type/save", self.Save)
	self.Router.GET("/movie_type/del", self.Delete)
}

type MovieTypeController struct {
	Router gin.IRouter
}

func (self *MovieTypeController) Get(c *gin.Context) {
	tagId := GetId(c)
	tag, err := movieTypeService.Get(tagId)
	if err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, tag)
}

func (self *MovieTypeController) Find(c *gin.Context) {
	page, limit := GetPageParams(c)
	typeKey := c.Query("keyword")
	tags, err := movieTypeService.FindByPage(page, limit, typeKey)
	if err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, tags)
}

func (self *MovieTypeController) Save(c *gin.Context) {
	req := sys.MovieType{}
	if err := BindJSON(c, &req); err != nil {
		ResponseErrorf(c, "请求的数据不合法: %s", err)
		return
	}

	if err := movieTypeService.Save(&req); err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}

func (self *MovieTypeController) Delete(c *gin.Context) {
	id := c.Query("id")
	Id, _ := strconv.ParseInt(id, 10, 64)
	if err := movieTypeService.Delete(Id); err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}
