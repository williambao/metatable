package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string
		Port int
	}

	IsProduction bool

	Database struct {
		IsMySQL      bool
		Host         string
		Port         string
		DatabaseName string
		Username     string
		Password     string
	}

	FileServer struct {
	}
}

func GetConfig() (*Config, error) {

	vp := viper.New()

	vp.SetDefault("server.host", "")
	vp.SetDefault("server.port", "9000")

	// 默认为sqlite数据库
	vp.SetDefault("database.is_mysql", false)
	vp.SetDefault("database.host", "./metatable.db")
	vp.SetDefault("database.port", 3306)

	vp.AutomaticEnv()

	vp.SetConfigName("config")           // name of config file (without extension)
	vp.SetConfigType("yaml")             // REQUIRED if the config file does not have the extension in the name
	vp.AddConfigPath("/etc/metatable/")  // path to look for the config file in
	vp.AddConfigPath("$HOME/.metatable") // call multiple times to add many search paths
	vp.AddConfigPath(".")                // optionally look for config in the working directory
	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Info("No config.yaml file found. use default config settings")
		} else {
			return nil, err
		}
	}

	config := Config{}
	config.Server.Host = vp.GetString("server.host")
	config.Server.Port = vp.GetInt("server.port")

	config.Database.IsMySQL = vp.GetBool("database.is_mysql")
	config.Database.Host = vp.GetString("database.host")
	config.Database.Port = vp.GetString("database.port")
	config.Database.DatabaseName = vp.GetString("database.database_name")
	config.Database.Username = vp.GetString("database.username")
	config.Database.Password = vp.GetString("database.password")

	return &config, nil
}

// func SetDefault(key string, value interface{}) {
// 	cfg.SetDefault(key, value)
// }
