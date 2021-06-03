/*
创建者：     Zero
创建时间：   2021/5/25
项目名称：   golang-study
*/
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var (
	err3          error
	eccPrivateKay *ecdsa.PrivateKey
	eccPublicKey  *ecdsa.PublicKey
)

//ecdsa 椭圆曲线密钥生成
func getEcdsaKey(keytype int) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	var err error
	var prk *ecdsa.PrivateKey
	var pub *ecdsa.PublicKey
	var curve elliptic.Curve //椭圆曲线
	switch keytype {
	case 1:
		curve = elliptic.P224()
	case 2:
		curve = elliptic.P256()
	case 3:
		curve = elliptic.P384()
	case 4:
		curve = elliptic.P521()
	default:
		err = errors.New("输入的签名key类型错误！key取值：\n 1：椭圆曲线224\n 2：椭圆曲线256\n 3：椭圆曲线384\n 4：椭圆曲线521\n")
		return nil, nil, err
	}
	prk, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	pub = &prk.PublicKey
	return prk, pub, err
}

//ecdsa椭圆曲线密钥认证中间件
func ecdsaAuthMiddleware() gin.HandlerFunc {
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
		claims, err := ecdsaParseToken(tokenString)
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

//解析ecdsa密钥
func ecdsaParseToken(tokenString string) (interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("无效的签名方法：%v", token.Method)
		}
		return eccPublicKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func checkToken3(c *gin.Context) {
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "响应成功！", "OK!"})
}

func getToken3(c *gin.Context) {
	var user User
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "参数错误！", err})
		return
	}
	token, err := ecdsaReleaseToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "生成Token错误！", err})
		return
	}
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "Token分发成功！", token})
}

//生成token
func ecdsaReleaseToken(user User) (interface{}, error) {
	claim := &MyClaims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),                         //发布时间
			Subject:   "User Token",                              //主题
			Issuer:    "Zero",                                    //发布者
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claim)
	signedString, err := token.SignedString(eccPrivateKay)
	return signedString, err
}
