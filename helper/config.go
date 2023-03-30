package helper

import "github.com/spf13/viper"

type Config struct {
	OpenAIKey      string `mapstructure:"OPEN_AI_KEY"`
	AppID          string `mapstructure:"APP_ID"`
	AppSecret      string `mapstructure:"APP_SECRET"`
	Token          string `mapstructure:"TOKEN"`
	EncodingAESKey string `mapstructure:"ENCODING_AES_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

var DefaultConfig *Config
