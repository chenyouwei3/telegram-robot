package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var Conf *Config

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/initialize/config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("err", err)
		panic(err)
	}
	err = viper.Unmarshal(&Conf)
	if err != nil {
		panic(err)
	}
}

type Config struct {
	Token string `yaml:"token"`
	Mysql MySQL  `yaml:"mysql"`
}

type MySQL struct {
	DriverName string `yaml:"driverName"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Database   string `yaml:"database"`
	UserName   string `yaml:"username"`
	Password   string `yaml:"password"`
}
