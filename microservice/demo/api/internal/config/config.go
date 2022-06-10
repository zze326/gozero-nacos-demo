package config

import "demo/common/config"
import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Mysql config.Mysql
}
