package sqlx

import (
	"fmt"
	"time"
)

type SqlConf struct {
	Name        string
	Host        string
	User        string
	Password    string
	Port        string
	Charset     string
	Driver      string
	Timeout     int
	MaxIdleConn int
	MaxOpenConn int
	MaxLifetime time.Duration
}

func (c SqlConf) NewSqlConn() SqlConn {

}

func (c SqlConf) GetConnectString() (connectString string) {
	result := ""
	switch cfg.Driver {
	case "mysql":
		timeout := cfg.Timeout
		if timeout <= 0 {
			timeout = 1
		}
		addr := l5parser.GetAddr(cfg.L5ID, cfg.Host+":"+cfg.Port)
		charset := cfg.Charset
		if charset == "" {
			charset = "utf8mb4"
		}
		// timeout 不是读的超时时间，是连接超时时间
		// * readTimeout 拍脑袋确定的 5s, 后续有更改可以配置文件追踪或其他方式修改。
		result = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=%ds&readTimeout=5s",
			cfg.User,
			cfg.Password,
			addr,
			cfg.Name,
			charset,
			timeout)
		break
	default:
		break
	}
	return result
}
