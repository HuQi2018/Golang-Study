/*
创建者：     Zero
创建时间：   2021/5/28
项目名称：   Zero-go
*/
package sys

import (
	"Zero-go/internal/model/dto"
	"Zero-go/internal/model/sys"
	"Zero-go/pkg/DB"
	"errors"
	"github.com/xormplus/xorm"
)

type MovieTypeService struct{}

var (
	SELECT_MOVIE_TYPE_TOTAL = "select_movie_type_total.stpl"
	SELECT_MOVIE_TYPE       = "select_movie_type.stpl"
)

func (self *MovieTypeService) Get(id int64) (*sys.MovieType, error) {
	tag := sys.MovieType{}
	has, err := DB.FindById(id, &tag)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("电影类型不存在")
	}

	return &tag, err
}

//分页模糊查询
func (self *MovieTypeService) FindByPage(page, limit int, keyword string) (pages *dto.Pages, err error) {
	params := map[string]interface{}{
		"offset":  GetOffset(page, limit),
		"limit":   limit,
		"keyword": LikeStr(keyword),
	}

	var total int64
	movie := make([]*sys.MovieType, 0)
	err = DB.PageBySqlTemplateClient(SELECT_MOVIE_TYPE, &params, &movie, SELECT_MOVIE_TYPE_TOTAL, &total)

	pages = &dto.Pages{Total: total, Data: movie}
	return pages, err
}

func (self *MovieTypeService) Save(tag *sys.MovieType) (err error) {
	session := NewSession()
	defer session.Close()

	return self.SaveTx(session, tag)
}

func (MovieTypeService) SaveTx(session *xorm.Session, tag *sys.MovieType) (err error) {
	count, err := DB.Where("name = ?", tag.Name).Count(&sys.MovieType{})
	if err != nil {
		return err
	}

	if tag.Id == 0 && count > 0 {
		return errors.New("电影类型信息已存在，请检查！")
	}

	if tag.Id == 0 { //修改电影信息
		_, err = DB.InsertOneTx(session, tag)
	} else {
		_, err = DB.UpdateByIdWithOmitTx(session, tag.Id, tag, "name")
	}

	return
}

func (self *MovieTypeService) Delete(id int64) (err error) {
	var count int64
	count, err = DB.Where("id = ?", id).Count(&sys.MovieType{})
	if err != nil {
		return
	}
	if count < 0 {
		return errors.New("电影类型信息不存在！")
	}

	movie := new(sys.MovieType)
	_, err = DB.DeleteById(id, movie)
	return
}
