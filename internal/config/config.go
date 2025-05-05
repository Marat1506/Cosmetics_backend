package config

import (
	"server/pkg/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"0.0.0.0"`
		Port   string `yaml:"port" env-default:"3000"`
	} `yaml:"listen"`
	MongoDB struct {
		Host             string `yaml:"host"`
		Port             string `yaml:"port"`
		Database         string `yaml:"database"`
		AuthDB           string `yaml:"auth_db"`
		Username         string `yaml:"username"`
		Password         string `yaml:"password"`
		Collection       string `yaml:"collection"`
		OrdersCollection string `yaml:"orders_collection"`
		Products         string `yaml:"products"`
	} `yaml:"mongodb"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
		logger.Infof("Full Config: %+v", instance)
	})
	return instance
}
