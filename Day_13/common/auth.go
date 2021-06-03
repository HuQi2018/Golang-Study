/*
创建者：     Zero
创建时间：   2021/5/25
项目名称：   golang-study
*/
package common

import (
	"github.com/dgrijalva/jwt-go"
	"golang-study/huqi/Day_13/model"
	"time"
)

var jwtKey = []byte("a_secret_key") //证书签名秘钥（该秘钥非常重要，如果client端有该秘钥，就可以签发证书了）

type MyClaims struct {
	UserId int
	jwt.StandardClaims
}

//解析ecdsa密钥
func ParseToken(tokenString string) (*jwt.Token, *MyClaims, error) {
	claims := &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}

//生成token  分发证书
func ReleaseToken(user model.UserBase) (string, error) {
	claim := &MyClaims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),                         //发布时间
			Subject:   "UserBase Token",                          //主题
			Issuer:    "Zero",                                    //发布者
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedString, err := token.SignedString(jwtKey)
	return signedString, err
}
