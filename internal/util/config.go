package util

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	TmpPath     string `mapstructure:"TMP_PATH"`
	Token       string `mapstructure:"TOKEN"`
	RESRGANPath string `mapstructure:"RESRGAN_PATH"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.BindEnv("TOKEN")
	viper.BindEnv("TMP_PATH")
	viper.BindEnv("RESRGAN_PATH")
	viper.AutomaticEnv()

	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.ReadInConfig()
	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}
	if config.Token == "" {
		err = fmt.Errorf("TOKEN is required but not set in environment variables or config file")
		return
	}
	return
}
