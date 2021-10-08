package main

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Host      string   `json:"host"      mapstructure:"host"`
	User      string   `json:"user"      mapstructure:"user"`
	Pwd       string   `json:"pwd"       mapstructure:"pwd"`
	DbName    string   `json:"db_name"   mapstructure:"dbname"`
	Collation string   `json:"collation" mapstructure:"collation"`
	Tables    []string `json:"tables"    mapstructure:"tables"`
}

var Cfg = Config{
	Host:      "",
	User:      "root",
	Pwd:       "",
	DbName:    "",
	Collation: "utf8mb4_unicode_ci",
	Tables:    nil,
}

func initConfig() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	viper.SetConfigFile(path + "/db.json")
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("未找到配置文件")
		} else {
			// Config file was found but another error was produced
			panic(err)
		}
	}

	if err = viper.Unmarshal(&Cfg); err != nil {
		panic(err)
	}
}
