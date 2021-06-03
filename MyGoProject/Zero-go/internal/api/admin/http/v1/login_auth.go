/**
 * Created by Wangwei on 2019-06-04 17:10.
 */

package v1

import (
	"Zero-go/internal/common"
	"Zero-go/internal/common/middleware/jwt"
	"Zero-go/internal/model/sys"
	service "Zero-go/internal/service/sys"
	"Zero-go/pkg/util/gosecurity"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/g/os/glog"
	"net/http"
)

var (
	myJwt       = &jwt.JWT{}
	userService = &service.UserService{}
)

type AuthController struct{}

func (AuthController) Login(c *gin.Context) {
	clientId := c.Request.Header.Get("X-Client-Id")
	clientId = common.FormatClientId(clientId)

	var loginInfoReq *sys.LoginInfoReq
	if BindJSON(c, &loginInfoReq) != nil {
		ResponseError(c, "无效的请求数据")
		return
	}

	password := gosecurity.MD5Password(loginInfoReq.Password)

	//glog.Println(password)
	user, err := userService.GetByAccount(loginInfoReq.Account, password)
	if err != nil {
		ResponseError(c, err)
		return
	}

	routeLinks, err := roleMenuService.FindRouteLinksByRole(user.RoleId, user.Id)
	if err != nil {
		glog.Error(err)
		ResponseError(c, err)
		return
	}

	tokenStr, exp, err := myJwt.GenUserToken(user.Id, user.Account,
		user.RoleId, user.OrgTypeId, user.OrgId, user.OrgName, clientId)
	if err != nil {
		ResponseError(c, "生成token失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        0,
		"user":        user,
		"route_links": routeLinks,
		"jwt": gin.H{
			"token":      tokenStr,
			"expires_in": exp,
		},
	})
}
