/*
创建者：     Zero
创建时间：   2021/5/25
项目名称：   golang-study
*/
package util

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	letters := []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
