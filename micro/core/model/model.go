package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/**
 * @Com www.github.com/sixstaredu
 * @Author 六星教育-shineyork老师
 */

func InitDb(config *Config) (DB *gorm.DB) {
	var err error

	username := config.Username
	password := config.Password
	host := config.Host
	port := config.Port
	dbname := config.Dbname

	dsn := fmt.Sprintf(DSN, username, password, host, port, dbname)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
	})
	if err != nil {
		fmt.Println("models/db.go:InitDb Fatal error mysql connect : ", err)
		return
	}
	initWuid(dsn)

	return
}
