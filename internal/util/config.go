package util

import "github.com/spf13/viper"

// All configuration of the application
type Config struct {
	TmpPath string `mapstructure:"TMP_PATH"`
	Token   string `mapstructure:"TOKEN"`
}

// Reads configuration from environment file then environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.BindEnv("TOKEN", "TMP_PATH")
	viper.AutomaticEnv()

	if err = viper.Unmarshal(&config); err != nil || config.Token == "" {
		viper.SetConfigName("app")
		viper.SetConfigType("env")
		viper.AddConfigPath(path)

		err = viper.ReadInConfig()
		if err != nil {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
