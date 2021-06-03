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
Token：令牌，是用户身份的验证方式。
	最简单的token组成：uid（用户唯一的身份标识）、time（当前时间的时间戳）、sign（签名）
token和session其实都是为了身份验证，session一本翻译为会话，而token更多的时候是翻译为令牌；
session服务器会保存一份，可能保存到缓存、文件、数据库；同样，session和token都是有过期时间一说，都需要去管理过期时间；
其实token与session的问题是一种时间与空间的博弈问题，session是空间换时间，而token是时间换空间。两者的选择要看具体情况而定。

虽然确实都是“客户端记录，每次访问携带”，但token很容易设计为自包含的，也就是说，后端不需要记录什么东西，每次一个无状态请求，
每次解密验证，每次当场得出合法/非法的结论。这一判断依据，出了固化在CS两端的一些逻辑之外，整个信息是自包含的。这才是真正的无状态。
而sessionid，一般都是一般随机字符串，而需要到后端去检索id的有效性。万一服务器重启导致内村里的session没了呢？

签名方法和key类型：
The HMAC signing method (HS256, HS384, HS512) //hash消息认证码
The RSA signing method (RS256, RS384, RS512)  //RSA非对称加密签名
The ECDSA signing method (ES256, ES384, ES512)  //椭圆曲线数字签名

文档地址：https://github.com/dgrijalva/jwt-go
*/

//消息认证签名实现token
type User struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}

type MyClaims struct {
	UserId string
	jwt.StandardClaims
}

type HttpResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var jwtKey = []byte("a_secret_key") //证书签名秘钥（该秘钥非常重要，如果client有该秘钥，就可以签发证书了。）

//hash消息认证中间件
func hamacAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := "Zero"
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, auth+":") { //验证Token不为空，并且以：Zero: 为前缀
			c.JSON(http.StatusUnauthorized, HttpResponse{http.StatusUnauthorized, "前缀错误，验证失败！", "Error!"})
			c.Abort()
			return
		}
		index := strings.Index(tokenString, auth+":") //找到token前缀对应的位置
		//真实token的值
		tokenString = tokenString[index+len(auth)+1:] //截取真实的Token（开始位置为：索引开始的位置+关键字符的长度+1（:的长度为1））
		token, claims, err := hamcParseToken(tokenString)
		if err != nil || !token.Valid { //解析错误或者过期等
			c.JSON(http.StatusUnauthorized, HttpResponse{http.StatusUnauthorized, "证书无效！", err})
			c.Abort()
			return
		}
		var hamcuser User
		err = c.Bind(&hamcuser)
		if err != nil {
			c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "请求数据参数有误！", err})
			c.Abort()
			return
		}
		if hamcuser.Id != claims.UserId {
			c.JSON(http.StatusUnauthorized, HttpResponse{http.StatusUnauthorized, "用户不存在！", err})
			c.Abort()
			return
		}
		c.Next()
	}
}

//解析Token
func hamcParseToken(tokenString string) (*jwt.Token, *MyClaims, error) {
	claims := &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}

//检查Token后的请求
func checkToken1(c *gin.Context) {
	c.JSON(http.StatusOK, "验证成功")
}

//获取Token
func getToken1(c *gin.Context) {
	var hamcuser User
	err := c.Bind(&hamcuser)
	fmt.Println(hamcuser.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpResponse{http.StatusBadRequest, "请求数据参数有误！", err})
		return
	}
	token, err := hamcReleaseToken(hamcuser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, HttpResponse{http.StatusOK, "token分发成功！", token})
}

//分发Token
func hamcReleaseToken(hamcuser User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) //截止时间：从当前时刻算起，7天
	claims := &MyClaims{
		UserId: hamcuser.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),     //发布时间
			Subject:   "User Token",          //主题
			Issuer:    "Zero",                //发布者
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //生成Token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
