package config

import (
	"sync"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

var (
	PATH = "./internal/config/config.yml"
	cfg  *Config
	once sync.Once
)

type Config struct {
	IsDebug bool `json:"is_debug"`
	Listen  struct {
		Type   string `json:"type"`
		BindIP string `json:"bind_ip"`
		Port   string `json:"port"`
	} `json:"listen"`
	Storage StorageConfig `json:"storage"`
}

type StorageConfig struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	Database    string `json:"database"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	MaxAttempts int    `json:"max_attempts"`
}

func GetConfig() *Config {
	logger := logging.GetLogger()
	once.Do(func() {
		cfg = &Config{}
		err := cleanenv.ReadConfig(PATH, cfg)
		if err != nil {
			help, _ := cleanenv.GetDescription(cfg, nil)
			logger.Info(LOG_INFO, help)
			logger.Fatal(LOG_ERROR, err)
		}
	})
	return cfg
}
