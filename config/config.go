package config

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"reflect"
	"scraper/internal/logger"
)

type Cfg struct {
	Port        int    `map:"PORT"`
	Host        string `map:"HOST"`
	ChromeSetup string `map:"CHROME_SETUP"`
}

var Config *Cfg

func setDefaults() {
	viper.SetDefault("PORT", 8080)
	viper.SetDefault("HOST", "0.0.0.0")
}

func GetConfig() *Cfg {
	ctx := context.Background()
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		typ := reflect.TypeOf(Config).Elem()
		for i := range typ.NumField() {
			err := viper.BindEnv(typ.Field(i).Tag.Get("map"))
			if err != nil {
				logger.ErrorCtx(ctx, "Failed to bind environment variable", logger.Field{Key: "error", Value: err})
				return nil
			}
		}
	}

	setDefaults()

	if err := viper.Unmarshal(&Config); err != nil {
		logger.ErrorCtx(ctx, "Failed to unmarshal configuration", logger.Field{Key: "error", Value: err})
		panic(err)
	}

	if errs := validator.New().Struct(Config); errs != nil {
		logger.ErrorCtx(ctx, "Invalid environment configuration", logger.Field{Key: "error", Value: errs})
		panic(errs)
	}
	return Config
}
