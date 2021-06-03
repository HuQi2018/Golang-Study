/*
创建者：     Zero
创建时间：   2021/5/25
项目名称：   golang-study
*/
package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

/*
	RAS密钥生成工具链接：http://www.metools.info/code/c80.html
*/

var (
	resPublicKey   []byte
	resPrivateKey  []byte
	err2_1, err2_2 error
)

func checkToken2(c *gin.Context) {
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "响应成功！", "OK!"})
}

//RSA验证中间件
//Zero:eyJhbGciOiJSUzI1NiIsIn
func rsaAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := "Zero"
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, auth+":") {
			c.JSON(http.StatusUnauthorized, HttpResponse{http.StatusUnauthorized, "无效的Token！", "Error!"})
			c.Abort()
			return
		}
		index := strings.Index(tokenString, auth+":")
		tokenString = tokenString[index+len(auth)+1:]
		claims, err := rsaParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, HttpResponse{http.StatusUnauthorized, "证书无效！", err})
			c.Abort()
			return
		}
		claimsValue := claims.(jwt.MapClaims) //断言
		rsauser := User{}
		err = c.Bind(&rsauser)
		if err != nil {
			c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "请求数据参数有误！", err})
			c.Abort()
			return
		}
		if claimsValue["UserId"] == nil || rsauser.Id != claimsValue["UserId"].(string) {
			c.JSON(http.StatusUnauthorized, HttpResponse{http.StatusUnauthorized, "用户不存在！", err})
			c.Abort()
			return
		}
		c.Next()
	}
}

//RSA Token验证
func rsaParseToken(tokenString string) (interface{}, error) {
	pem, err := jwt.ParseRSAPublicKeyFromPEM(resPublicKey)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		if _, OK := token.Method.(*jwt.SigningMethodRSA); !OK {
			return nil, fmt.Errorf("解析的方法错误")
		}
		return pem, err
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func getToken2(c *gin.Context) {
	rsaUser := User{}
	err := c.Bind(&rsaUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "参数错误！", err})
		return
	}
	token, err := resaReleaseToken(rsaUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "生成Token错误！", err})
		return
	}
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "Token分发成功！", token})
}

//分发Token
func resaReleaseToken(user User) (interface{}, error) {
	tokenGen, err := rasJwtTokenGen(user.Id)
	return tokenGen, err
}

//生成Token
func rasJwtTokenGen(id string) (interface{}, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(resPrivateKey)
	if err != nil {
		return nil, err
	}
	claim := &MyClaims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),                         //发布时间
			Subject:   "User Token",                              //主题
			Issuer:    "Zero",                                    //发布者
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	signedString, err := token.SignedString(privateKey)
	return signedString, err
}
