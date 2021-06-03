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
	movieTagService = &sys_service.MovieTagService{}
)

func (self *MovieTagController) Setup() {
	self.Router.GET("/movie_tag/get", self.Get)
	self.Router.GET("/movie_tag/query", self.Find)
	self.Router.POST("/movie_tag/save", self.Save)
	self.Router.GET("/movie_tag/del", self.Delete)
}

type MovieTagController struct {
	Router gin.IRouter
}

func (self *MovieTagController) Get(c *gin.Context) {
	tagId := GetId(c)
	tag, err := movieTagService.Get(tagId)
	if err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, tag)
}

func (self *MovieTagController) Find(c *gin.Context) {
	page, limit := GetPageParams(c)
	tagKey := c.Query("keyword")
	tags, err := movieTagService.FindByPage(page, limit, tagKey)
	if err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, tags)
}

func (self *MovieTagController) Save(c *gin.Context) {
	req := sys.MovieTag{}
	if err := BindJSON(c, &req); err != nil {
		ResponseErrorf(c, "请求的数据不合法: %s", err)
		return
	}

	if err := movieTagService.Save(&req); err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}

func (self *MovieTagController) Delete(c *gin.Context) {
	id := c.Query("id")
	Id, _ := strconv.ParseInt(id, 10, 64)
	if err := movieTagService.Delete(Id); err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}
