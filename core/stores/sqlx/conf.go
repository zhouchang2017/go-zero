package sqlx

import (
	"fmt"
	"time"
)

type SqlConf struct {
	Name        string `validate:"required"`
	Addr        string `validate:"required"`
	User        string `validate:"required"`
	Password    string
	Charset     string        `default:"utf8mb4"`
	Driver      string        `default:"mysql"`
	Timeout     time.Duration `default:"1s"`
	ReadTimeout time.Duration `default:"5s"`
	MaxIdleConn int           `default:"64"`
	MaxOpenConn int           `default:"64"`
	MaxLifetime time.Duration `default:"1m"`
}

func (c *SqlConf) sourceName() string {
	return fmt.Sprintf("%s://%s@%s/%s",
		c.Driver,
		c.User,
		c.Addr,
		c.Name)
}
