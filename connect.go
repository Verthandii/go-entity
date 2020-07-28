package main

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Connection struct {
	DB *gorm.DB
}

func Connect() Connection {
	cfg := mysql.Config{
		Addr:                 Cfg.Host,
		User:                 Cfg.User,
		Passwd:               Cfg.Pwd,
		DBName:               Cfg.DbName,
		Net:                  "tcp",
		Collation:            Cfg.Collation,
		AllowNativePasswords: true,
	}
	db, err := gorm.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}
	return Connection{DB: db}
}

func (db *Connection) Close() {
	db.DB.Close()
}
