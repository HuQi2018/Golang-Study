package main

import (
	"MyGoProject/common"
	"MyGoProject/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
)

func init() {

}

func main() {
	InitConfig()    //配置项的加载
	common.InitDB() //数据库初始化（只初始化一次）

	r := gin.Default()

	//HMAC Token分发
	r.POST("/getToken1", getToken1)

	//HMAC Token认证
	r.POST("/checkToken1", hamacAuthMiddleware(), checkToken1)

	//RSA Token分发
	r.POST("/getToken2", getToken2)

	//RSA Token认证
	r.POST("/checkToken2", rsaAuthMiddleware(), checkToken2)

	//ECDSA Token分发
	r.POST("/getToken3", getToken3)

	//ECDSA Token认证
	r.POST("/checkToken3", ecdsaAuthMiddleware(), checkToken3)

	r = router.CollectRouter(r)
	var err error
	port := viper.GetString("server.port")
	if port != "" {
		err = r.Run(":" + port)
	} else {
		err = r.Run() //默认8080
	}
	if err != nil {
		fmt.Println("服务进程启动失败！")
	}
}

func InitConfig() {
	workDir, _ := os.Getwd()                  //获取目录对应的路径
	viper.SetConfigName("application")        //配置文件名
	viper.SetConfigType("yml")                //配置文件类型
	viper.AddConfigPath("huqi/Day_13/config") //执行go run的对应路径配置
	fmt.Println(workDir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
