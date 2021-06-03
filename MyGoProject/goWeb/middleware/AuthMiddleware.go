/*
创建者：     Zero
创建时间：   2021/5/26
项目名称：   goWeb
*/
package middleware

import (
	"MyGoProject/common"
	"MyGoProject/global"
	"MyGoProject/model"
	"github.com/gin-gonic/gin"
	"strings"
)

//密钥认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := "Zero"
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, auth+":") {
			global.Fail(c, "Error!", "无效的Token！")
			//c.JSON(http.StatusUnauthorized, HttpResponse{http.StatusUnauthorized, "无效的Token！", "Error!"})
			c.Abort()
			return
		}
		index := strings.Index(tokenString, auth+":")
		tokenString = tokenString[index+len(auth)+1:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid { //解析错误或者过期等
			global.Fail(c, err, "证书无效！")
			//c.JSON(http.StatusUnauthorized, HttpResponse{http.StatusUnauthorized, "证书无效！", err})
			c.Abort()
			return
		}

		//验证通过后获取claims中的userId
		userId := claims.UserId
		//判定
		var user model.UserInfo
		common.MyDB.First(&user, userId)
		if user.ID == 0 {
			global.Fail(c, err, "用户不存在！")
			//c.JSON(http.StatusUnauthorized, HttpResponse{http.StatusUnauthorized, "用户不存在！", err})
			c.Abort()
			return
		}
		c.Set("user", user) //将key-value值存储到context中
		c.Next()
	}
}
