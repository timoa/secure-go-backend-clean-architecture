package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPass                 string `mapstructure:"DB_PASS"`
	DBName                 string `mapstructure:"DB_NAME"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	_ = viper.BindEnv("APP_ENV")
	_ = viper.BindEnv("SERVER_ADDRESS")
	_ = viper.BindEnv("CONTEXT_TIMEOUT")
	_ = viper.BindEnv("DB_HOST")
	_ = viper.BindEnv("DB_PORT")
	_ = viper.BindEnv("DB_USER")
	_ = viper.BindEnv("DB_PASS")
	_ = viper.BindEnv("DB_NAME")
	_ = viper.BindEnv("ACCESS_TOKEN_EXPIRY_HOUR")
	_ = viper.BindEnv("REFRESH_TOKEN_EXPIRY_HOUR")
	_ = viper.BindEnv("ACCESS_TOKEN_SECRET")
	_ = viper.BindEnv("REFRESH_TOKEN_SECRET")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("No .env file found, using environment variables only: %v", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
