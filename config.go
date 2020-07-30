package main

import "github.com/spf13/viper"

type Config struct {
	Host      string
	User      string
	Pwd       string
	DbName    string
	Collation string
	Tables    []string
}

var Cfg = Config{
	Host:      "",
	User:      "root",
	Pwd:       "",
	DbName:    "",
	Collation: "utf8mb4_unicode_ci",
	Tables:    nil,
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
	if user := parser.GetString("user"); user != "" {
		Cfg.User = user
	}
	Cfg.Pwd = parser.GetString("pwd")
	Cfg.DbName = parser.GetString("dbname")
	if Cfg.DbName == "" {
		panic("DbName must be specified")
	}
	if collation := parser.GetString("collation"); collation != "" {
		Cfg.Collation = collation
	}
	Cfg.Tables = parser.GetStringSlice("tables")
}
