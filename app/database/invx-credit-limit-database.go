package database

import (
	"creditlimit-connector/app/configs"
	"fmt"
	"sync"
	"time"

	"creditlimit-connector/app/log"

	mysqlFormat "github.com/go-sql-driver/mysql"
	mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	onceCreditLimitDB sync.Once
	gormCreditLimit   *gorm.DB
)

func InitINVXCreditLimitDatabase() *gorm.DB {
	onceCreditLimitDB.Do(func() { // <-- atomic, does not allow repeating
		conf := configs.Conf.Database
		mySqlConf := mysqlFormat.Config{
			User:                 conf.INVXCreditLimitDatabase.User,
			Passwd:               conf.INVXCreditLimitDatabase.Password,
			Net:                  "tcp",
			Addr:                 fmt.Sprintf("%s:%d", conf.INVXCreditLimitDatabase.Host, conf.INVXCreditLimitDatabase.Port),
			DBName:               conf.INVXCreditLimitDatabase.Name,
			AllowNativePasswords: true,
			ParseTime:            true,
			TLSConfig:            configs.Conf.Database.SSLMode,
		}

		dsn := mySqlConf.FormatDSN()
		gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to open to the database: %v", err)
		}

		log.Info("Connect credit limit database success")

		gormCreditLimit = gormDB

		rawDB, err := gormDB.DB()
		if err != nil {
			log.Fatalf("failed to connect credit limit database: %v", err)
		}
		// check connection
		rawDB.Ping()
		// set connection pool
		rawDB.SetMaxIdleConns(conf.DBMaxIdleConnections)
		rawDB.SetMaxOpenConns(conf.DBMaxConnections)
		rawDB.SetConnMaxLifetime(time.Duration(conf.DBMaxLifetimeConnections))

	})
	return gormCreditLimit

}
