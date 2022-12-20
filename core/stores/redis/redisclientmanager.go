package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"

	red "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/syncx"
)

const (
	defaultDatabase = 0
	maxRetries      = 3
	idleConns       = 8
)

var clientManager = syncx.NewResourceManager()
var Dialer func(ctx context.Context, network, addr string) (net.Conn, error)

func getClient(r *Redis) (*red.Client, error) {
	val, err := clientManager.GetResource(fmt.Sprintf("%s_%d", r.Addr, r.DB), func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if r.tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		store := red.NewClient(&red.Options{
			Addr:         r.Addr,
			Password:     r.Pass,
			DB:           r.DB,
			MaxRetries:   maxRetries,
			MinIdleConns: idleConns,
			TLSConfig:    tlsConfig,
			Dialer:       Dialer,
		})
		store.AddHook(durationHook)

		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.Client), nil
}
