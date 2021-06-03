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
	movieService = &sys_service.MovieService{}
)

func (self *MovieController) Setup() {
	self.Router.GET("/movie/query", self.Find)
	self.Router.POST("/movie/save", self.Save)
	self.Router.GET("/movie/get", self.Get)
	self.Router.GET("/movie/del", self.Delete)
}

type MovieController struct {
	Router gin.IRouter
}

func (*MovieController) Find(c *gin.Context) {
	page, limit := GetPageParams(c)
	movieId := c.Query("keyword")
	movies, err := movieService.FindByPage(page, limit, movieId)
	if err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, movies)
}

func (self *MovieController) Get(c *gin.Context) {
	movieId := GetId(c)
	movie, err := movieService.FindById(movieId)
	if err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, movie)
}

func (*MovieController) Save(c *gin.Context) {
	req := sys.MovieInfo{}
	if err := BindJSON(c, &req); err != nil {
		ResponseErrorf(c, "请求的数据不合法: %s", err)
		return
	}

	if err := movieService.Save(&req); err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}

func (*MovieController) Delete(c *gin.Context) {

	id := c.Query("id")
	Id, _ := strconv.ParseInt(id, 10, 64)
	if err := movieService.Delete(Id); err != nil {
		ResponseError(c, err)
		return
	}

	ResponseOK(c, "")
}
