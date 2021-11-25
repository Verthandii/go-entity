package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connectMySQL() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Cfg.Username, Cfg.Password, Cfg.Host, Cfg.Port, Cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rawDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	if err = rawDB.Ping(); err != nil {
		panic(err)
	}
	return db
}
