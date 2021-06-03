/**
 * Created by Wangwei on 2019-05-30 18:10.
 */

package conf

import (
	"Zero-go/pkg/DB"
	"Zero-go/pkg/util/gopath"
	"flag"
	"fmt"

	"github.com/gogf/gf/g/os/glog"

	"github.com/BurntSushi/toml"
)

var (
	confPath   string
	AppConf    = AppConfig{}
	ServerConf = ServerConfig{}
	MySQLConf  = DB.MySQLConfig{}
	LogConf    = LogConfig{}
)

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

// Init conf.
func Init() {
	if confPath == "" {
		confPath = gopath.FindParentPath("configs", "application.toml")
	}

	appPath := fmt.Sprintf("%s/application.toml", confPath)
	httpPath := fmt.Sprintf("%s/server.toml", confPath)
	mysqlPath := fmt.Sprintf("%s/mysql.toml", confPath)
	logPath := fmt.Sprintf("%s/log.toml", confPath)

	tomlDecodeFile(appPath, &AppConf)
	tomlDecodeFile(httpPath, &ServerConf)
	tomlDecodeFile(mysqlPath, &MySQLConf)
	tomlDecodeFile(logPath, &LogConf)
}

func tomlDecodeFile(fpath string, v interface{}) {
	_, err := toml.DecodeFile(fpath, v)
	if err != nil {
		glog.Fatalf("err: %v\n", err)
		panic(err)
	}
}
