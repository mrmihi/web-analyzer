package config

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"scraper/internal/logger"
)

type Cfg struct {
	Port        string `mapstructure:"PORT"`
	Host        string `mapstructure:"HOST"`
	ChromeSetup string `mapstructure:"CHROME_SETUP"`
	Leakless    bool   `mapstructure:"LEAKLESS"`
}

var Config *Cfg

func setDefaults() {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("LEAKLESS", true)
}

func GetConfig() *Cfg {
	ctx := context.Background()

	viper.AutomaticEnv()
	bindEnvVariables()

	setDefaults()

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	_ = viper.ReadInConfig()

	Config = &Cfg{}
	if err := viper.Unmarshal(Config); err != nil {
		logger.ErrorCtx(ctx, "Failed to unmarshal configuration", logger.Field{Key: "error", Value: err})
		panic(err)
	}

	if err := validator.New().Struct(Config); err != nil {
		logger.ErrorCtx(ctx, "Invalid environment configuration", logger.Field{Key: "error", Value: err})
		panic(err)
	}

	return Config
}

func bindEnvVariables() {
	_ = viper.BindEnv("PORT")
	_ = viper.BindEnv("HOST")
	_ = viper.BindEnv("CHROME_SETUP")
	_ = viper.BindEnv("LEAKLESS")
}
