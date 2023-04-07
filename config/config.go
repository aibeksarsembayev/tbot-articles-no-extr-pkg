package config

import "github.com/spf13/viper"

// Config
type Config struct {
	Debug    bool     `mapstructure:"debug"`
	Server   Server   `mapstructure:"server"`
	Context  Context  `mapstructure:"context"`
	Database Database `mapstructure:"database"`
	DBTest   DBTest   `mapstructure:"db_test"`
	TgBot    TgBot    `mapstructure:"tgbot"`
}

type Server struct {
	Address string `mapstructure:"address"`
}

type Context struct {
	Timeout int `mapstructure:"timeout"`
}

type Database struct {
	DBHost string `mapstructure:"host"`
	DBPort string `mapstructure:"port"`
	DBUser string `mapstructure:"user"`
	DBPass string `mapstructure:"pass"`
	DBName string `mapstructure:"name"`
}

type DBTest struct {
	DBHost string `mapstructure:"host"`
	DBPort string `mapstructure:"port"`
	DBUser string `mapstructure:"user"`
	DBPass string `mapstructure:"pass"`
	DBName string `mapstructure:"name"`
}

type TgBot struct {
	Host           string `mapstructure:"host"`
	BatchSize      int    `mapstructure:"batch_size"`
	Token          string `mapstructure:"token"`
	DigestChatID   int    `mapstructure:"digest_chat_id"`
	APIParsePeriod int    `mapstructure:"api_parse_period"`
}

var vp *viper.Viper

// Load Config from JSON into stucture ...
func LoadConfig() (Config, error) {
	vp = viper.New()

	var config Config

	vp.SetConfigName("config")
	vp.SetConfigType("json")
	vp.AddConfigPath(".")
	vp.AddConfigPath("./config")
	vp.AddConfigPath("../../config") // to test subfolders

	err := vp.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}
	return config, err
}
