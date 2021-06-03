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

type MovieService struct{}

var (
	SELECT_MOVIE_TOTAL = "select_movie_total.stpl"
	SELECT_MOVIE       = "select_movie.stpl"
)

func (self *MovieService) Get(id int64) (*sys.MovieInfo, error) {
	movie := sys.MovieInfo{}
	has, err := DB.FindById(id, &movie)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("电影不存在")
	}

	return &movie, err
}

//分页模糊查询
func (self *MovieService) FindByPage(page, limit int, keyword string) (pages *dto.Pages, err error) {
	params := map[string]interface{}{
		"offset":  GetOffset(page, limit),
		"limit":   limit,
		"keyword": LikeStr(keyword),
	}

	var total int64
	movie := make([]*sys.MovieInfo, 0)
	err = DB.PageBySqlTemplateClient(SELECT_MOVIE, &params, &movie, SELECT_MOVIE_TOTAL, &total)

	pages = &dto.Pages{Total: total, Data: movie}
	return pages, err
}

//通过movieId查询
func (self *MovieService) FindById(id int64) (*sys.MovieInfo, error) {
	movie := new(sys.MovieInfo)
	ok, err := DB.FindById(id, movie)
	if !ok {
		return nil, errors.New("电影ID:" + string(movie.Id) + "不存在")
	}

	return movie, err
}

func (self *MovieService) Save(movie *sys.MovieInfo) (err error) {
	session := NewSession()
	defer session.Close()

	return self.SaveTx(session, movie)
}

func (MovieService) SaveTx(session *xorm.Session, movie *sys.MovieInfo) (err error) {
	count, err := DB.Where("movie_id = ?", movie.MovieId).Count(&sys.MovieInfo{})
	if err != nil {
		return err
	}

	if movie.Id == 0 && count > 0 {
		return errors.New("电影信息已存在，请检查！")
	}

	if movie.Id == 0 { //修改电影信息
		_, err = DB.InsertOneTx(session, movie)
	} else {
		_, err = DB.UpdateByIdWithOmitTx(session, movie.Id, movie, "movie_id")
	}

	return
}

func (self *MovieService) Delete(id int64) (err error) {
	var count int64
	count, err = DB.Where("id = ?", id).Count(&sys.MovieInfo{})
	if err != nil {
		return
	}
	if count < 0 {
		return errors.New("电影信息不存在！")
	}

	movie := new(sys.MovieInfo)
	_, err = DB.DeleteById(id, movie)
	return
}

//func (self *MovieService) ResetPassword(account string) (err error) {
//
//	user := new(sys.MovieInfo)
//	ok, err := DB.Where("`account`=? ", account).FindOne(user)
//	if err != nil {
//		return err
//	}
//	if !ok {
//		return errors.New("无此用户")
//	}
//
//	user.Password = gosecurity.MD5Password(conf.AppConf.DefaultPassword)
//	_, err = DB.UpdateById(user.Id, user)
//
//	return err
//}
//
//func (self *MovieService) UpdatePassword(req *sys.PasswordReq) (err error) {
//
//	user := new(sys.MovieInfo)
//	ok, err := DB.Where("`account`=? and `password`=?", req.Account, gosecurity.MD5Password(req.OldPassword)).FindOne(user)
//	if err != nil {
//		return err
//	}
//	if !ok {
//		return errors.New("密码错误")
//	}
//
//	user.Password = gosecurity.MD5Password(req.Password)
//	_, err = DB.UpdateById(user.Id, user)
//
//	return err
//}
