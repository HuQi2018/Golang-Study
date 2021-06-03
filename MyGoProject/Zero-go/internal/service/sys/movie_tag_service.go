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

type MovieTagService struct{}

var (
	SELECT_MOVIE_TAG_TOTAL = "select_movie_tag_total.stpl"
	SELECT_MOVIE_TAG       = "select_movie_tag.stpl"
)

func (self *MovieTagService) Get(id int64) (*sys.MovieTag, error) {
	tag := sys.MovieTag{}
	has, err := DB.FindById(id, &tag)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("标签不存在")
	}

	return &tag, err
}

//分页模糊查询
func (self *MovieTagService) FindByPage(page, limit int, keyword string) (pages *dto.Pages, err error) {
	params := map[string]interface{}{
		"offset":  GetOffset(page, limit),
		"limit":   limit,
		"keyword": LikeStr(keyword),
	}

	var total int64
	tag := make([]*sys.MovieTag, 0)
	err = DB.PageBySqlTemplateClient(SELECT_MOVIE_TAG, &params, &tag, SELECT_MOVIE_TAG_TOTAL, &total)

	pages = &dto.Pages{Total: total, Data: tag}
	return pages, err
}

func (self *MovieTagService) Save(tag *sys.MovieTag) (err error) {
	session := NewSession()
	defer session.Close()

	return self.SaveTx(session, tag)
}

func (MovieTagService) SaveTx(session *xorm.Session, tag *sys.MovieTag) (err error) {
	count, err := DB.Where("name = ?", tag.Name).Count(&sys.MovieTag{})
	if err != nil {
		return err
	}

	if tag.Id == 0 && count > 0 {
		return errors.New("标签信息已存在，请检查！")
	}

	if tag.Id == 0 { //修改电影信息
		_, err = DB.InsertOneTx(session, tag)
	} else {
		_, err = DB.UpdateByIdWithOmitTx(session, tag.Id, tag, "id")
	}

	return
}

func (self *MovieTagService) Delete(id int64) (err error) {
	var count int64
	count, err = DB.Where("id = ?", id).Count(&sys.MovieTag{})
	if err != nil {
		return
	}
	if count < 0 {
		return errors.New("标签信息不存在！")
	}

	movie := new(sys.MovieTag)
	_, err = DB.DeleteById(id, movie)
	return
}
