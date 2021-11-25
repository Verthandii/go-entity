package main

import (
	"io/fs"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Host     string   `mapstructure:"host"`
	Port     int      `mapstructure:"port"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
	DBName   string   `mapstructure:"dbname"`
	Tables   []string `mapstructure:"tables"`
	Output   string   `mapstructure:"output"`
	Terminal bool     `mapstructure:"terminal"`
}

var Cfg = Config{
	Host:     "127.0.0.1",
	Port:     3306,
	Username: "root",
	Password: "root",
	DBName:   "test",
	Tables:   nil,
	Output:   "./model",
	Terminal: false,
}

func initConfig() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	viper.SetConfigFile(path + "/db.json")
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(*fs.PathError); ok {
			panic("未找到配置文件")
		} else {
			panic(err)
		}
	}

	if err = viper.Unmarshal(&Cfg); err != nil {
		panic(err)
	}

	if Cfg.Output[len(Cfg.Output)-1] != '/' {
		Cfg.Output += "/"
	}
}
