package initializers

import "github.com/spf13/viper"

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	StaticPath     string `mapstructure:"STATIC_PATH"`
	TestDBHost     string `mapstructure:"TEST_POSTGRES_HOST"`
	TestDBPort     string `mapstructure:"TEST_POSTGRES_PORT"`
	TestDBName     string `mapstructure:"TEST_POSTGRES_DB"`
	Dev            bool   `mapstructure:"DEV"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
