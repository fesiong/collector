package services

import (
	"collector/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

func initDB() {
	fmt.Println(config.MySQLConfig.Url)
	db, err := gorm.Open("mysql", config.MySQLConfig.Url)
	if err != nil {
		fmt.Println( config.MySQLConfig, err.Error())
		os.Exit(-1)
	}

	if config.ServerConfig.Env == "development" {
		db.LogMode(true)
	}
	db.DB().SetMaxIdleConns(config.MySQLConfig.MaxIdleConnections)
	db.DB().SetMaxOpenConns(config.MySQLConfig.MaxOpenConnections)
	db.DB().SetConnMaxLifetime(-1) //不重新利用，可以执行得更快

	//禁用复数表名
	db.SingularTable(true)

	//统一加前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return config.MySQLConfig.TablePrefix + defaultTableName
	}

	DB = db
}

var DB * gorm.DB

func init() {
	initDB()
}
