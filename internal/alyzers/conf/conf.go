package conf

import (
	"fmt"
	"github.com/alyzers/alyzers/pkg/http"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"sync"

	"github.com/alyzers/alyzers/pkg/cache"
	"github.com/alyzers/alyzers/pkg/database"
	"github.com/alyzers/alyzers/pkg/log"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/9/8 23:20
 * @file: conf.go
 * @description:
 */

type AppConfig struct {
	Log      log.LogConfig
	Http     http.Http
	Database database.Database
	Redis    cache.Redis
}

var (
	cfg  AppConfig
	once sync.Once
)

func NewConf(confDir string) AppConfig {
	once.Do(func() {
		var err error
		cfg, err = LoadConfigFile(confDir)
		if err != nil {
			panic(fmt.Sprintf("load conf file error: %s", err))
		}
	})
	return cfg
}

// LoadConfigFile load conf file
func LoadConfigFile(confDir string) (AppConfig, error) {

	config := viper.New()
	config.SetConfigFile(confDir) //文件名
	config.AddConfigPath("./conf.d")
	config.SetConfigName("config")
	config.SetConfigType("toml")
	if err := config.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("failed to read configuration file: %v", err)
	}

	config.WatchConfig()
	config.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("The configuration changes, re -analyze the configuration file: %s", e.Name)
		if err := config.Unmarshal(&cfg); err != nil {
			_ = fmt.Errorf("failed to unmarshal configuration file: %v", err)
		}
	})
	if err := config.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal configuration file: %v", err)
	}
	fmt.Printf("[Init] config file path: %s\n", confDir)

	return cfg, nil
}

func GetString(key string) string {
	return viper.GetString(key)
}
