/*
创建者：     Zero
创建时间：   2021/5/25
项目名称：   golang-study
*/
package middleware

import (
	"github.com/gin-gonic/gin"
	"golang-study/huqi/Day_13/common"
	"golang-study/huqi/Day_13/model"
	"golang-study/huqi/Day_13/response"
	"strings"
)

//ecdsa椭圆曲线密钥认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := "Zero"
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, auth+":") {
			response.Fail(c, "Error!", "无效的Token！")
			//c.JSON(http.StatusUnauthorized, HttpResponse{http.StatusUnauthorized, "无效的Token！", "Error!"})
			c.Abort()
			return
		}
		index := strings.Index(tokenString, auth+":")
		tokenString = tokenString[index+len(auth)+1:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid { //解析错误或者过期等
			response.Fail(c, err, "证书无效！")
			//c.JSON(http.StatusUnauthorized, HttpResponse{http.StatusUnauthorized, "证书无效！", err})
			c.Abort()
			return
		}

		//验证通过后获取claims中的userId
		userId := claims.UserId
		//判定
		var user model.UserBase
		common.MyDB.First(&user, userId)
		if user.ID == 0 {
			response.Fail(c, err, "用户不存在！")
			//c.JSON(http.StatusUnauthorized, HttpResponse{http.StatusUnauthorized, "用户不存在！", err})
			c.Abort()
			return
		}
		c.Set("user", user) //将key-value值存储到context中
		c.Next()
	}
}
