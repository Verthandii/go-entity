package main

import "github.com/spf13/viper"

type Config struct {
	Host   string
	User   string
	Pwd    string
	DbName string
	Tables []string
}

var Cfg = Config{
	Host:   "127.0.0.1",
	User:   "root",
	Pwd:    "123123",
	DbName: "test",
	Tables: nil,
}

func init() {
	parser := viper.New()
	parser.SetConfigName("db")
	parser.SetConfigType("json")
	parser.AddConfigPath(".")
	if err := parser.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("未找到配置文件")
		} else {
			// Config file was found but another error was produced
			panic(err)
		}
	}

	Cfg.Host = parser.GetString("host")
	Cfg.User = parser.GetString("user")
	Cfg.Pwd = parser.GetString("pwd")
	Cfg.DbName = parser.GetString("dbname")
	Cfg.Tables = parser.GetStringSlice("tables")
}
