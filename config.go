package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Name string
}

// 读取配置
func (c *Config) InitConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")

	return viper.ReadInConfig()
}

// 监控配置改动
func (c *Config) WatchConfig(change chan int) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file is changed: %s", e.Name)
		change <- 1
	})
}