package database

import (
	"creditlimit-connector/app/configs"
	"fmt"
	"net/url"
	"sync"
	"time"

	"creditlimit-connector/app/log"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	onceINVXCenterDB sync.Once
	gormINVXCenter   *gorm.DB
)

func InitINVXCenterDatabase() *gorm.DB {

	onceINVXCenterDB.Do(func() { // <-- atomic, does not allow repeating
		conf := configs.Conf.Database
		invxConf := configs.Conf.Database.INVXCenterDatabase

		u := &url.URL{
			Scheme:     "sqlserver",
			User:       url.UserPassword(invxConf.User, invxConf.Password),
			Host:       fmt.Sprintf("%s:%d", invxConf.Host, invxConf.Port),
			ForceQuery: true,
			RawQuery:   "database=" + invxConf.Name,
		}

		dsn := u.String()
		gormDB, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})
		if err != nil {
			log.Fatalf("failed to open to the database: %v", err)
		}

		log.Info("Connect INVX center database success")

		gormINVXCenter = gormDB

		rawDB, err := gormDB.DB()
		if err != nil {
			log.Fatalf("failed to connect to INVX center database: %v", err)
		}
		// check connection
		rawDB.Ping()
		// set connection pool
		rawDB.SetMaxIdleConns(conf.DBMaxIdleConnections)
		rawDB.SetMaxOpenConns(conf.DBMaxConnections)
		rawDB.SetConnMaxLifetime(time.Duration(conf.DBMaxLifetimeConnections))

	})
	return gormINVXCenter

}
