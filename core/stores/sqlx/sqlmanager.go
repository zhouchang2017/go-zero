package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/named"
	"io"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/syncx"
)

const (
	maxIdleConns = 64
	maxOpenConns = 64
	maxLifetime  = time.Minute
)

var connManager = syncx.NewResourceManager()

type pingedDB struct {
	*sql.DB
	once sync.Once
}

func getCachedSqlConnFromConf(conf SqlConf) (*pingedDB, error) {
	val, err := connManager.GetResource(conf.sourceName(), func() (io.Closer, error) {
		conn, err := newDBConnectionFromConf(conf)
		if err != nil {
			return nil, err
		}

		return &pingedDB{
			DB: conn,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*pingedDB), nil
}

func getCachedSqlConn(driverName, server string) (*pingedDB, error) {
	val, err := connManager.GetResource(server, func() (io.Closer, error) {
		conn, err := newDBConnection(driverName, server)
		if err != nil {
			return nil, err
		}

		return &pingedDB{
			DB: conn,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*pingedDB), nil
}

func getSqlConnFromConf(conf SqlConf) (*sql.DB, error) {
	pdb, err := getCachedSqlConnFromConf(conf)
	if err != nil {
		return nil, err
	}

	pdb.once.Do(func() {
		err = pdb.Ping()
	})
	if err != nil {
		return nil, err
	}

	return pdb.DB, nil
}

func getSqlConn(driverName, server string) (*sql.DB, error) {
	pdb, err := getCachedSqlConn(driverName, server)
	if err != nil {
		return nil, err
	}

	pdb.once.Do(func() {
		err = pdb.Ping()
	})
	if err != nil {
		return nil, err
	}

	return pdb.DB, nil
}

func newDBConnection(driverName, datasource string) (*sql.DB, error) {
	conn, err := sql.Open(driverName, datasource)
	if err != nil {
		return nil, err
	}

	// we need to do this until the issue https://github.com/golang/go/issues/9851 get fixed
	// discussed here https://github.com/go-sql-driver/mysql/issues/257
	// if the discussed SetMaxIdleTimeout methods added, we'll change this behavior
	// 8 means we can't have more than 8 goroutines to concurrently access the same database.
	conn.SetMaxIdleConns(maxIdleConns)
	conn.SetMaxOpenConns(maxOpenConns)
	conn.SetConnMaxLifetime(maxLifetime)

	return conn, nil
}

func newDBConnectionFromConf(conf SqlConf) (*sql.DB, error) {
	instance, err := named.GetGlobalResolver().GetInstance(context.Background(), conf.Addr)
	if err != nil {
		return nil, err
	}

	datasource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=%s&readTimeout=%s",
		conf.User,
		conf.Password,
		instance.GetEndpoint(),
		conf.Name,
		conf.Charset,
		conf.Timeout, conf.ReadTimeout)

	now := time.Now()
	conn, err := sql.Open(conf.Driver, datasource)
	cost := time.Now().Sub(now)
	instance.CallResultReport(err, cost)
	if err != nil {
		return nil, err
	}

	// we need to do this until the issue https://github.com/golang/go/issues/9851 get fixed
	// discussed here https://github.com/go-sql-driver/mysql/issues/257
	// if the discussed SetMaxIdleTimeout methods added, we'll change this behavior
	// 8 means we can't have more than 8 goroutines to concurrently access the same database.
	conn.SetMaxIdleConns(conf.MaxIdleConn)
	conn.SetMaxOpenConns(conf.MaxOpenConn)
	conn.SetConnMaxLifetime(conf.MaxLifetime)

	return conn, nil
}
