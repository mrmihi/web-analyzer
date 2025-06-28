package config

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
	"scraper/internal/logger"
	"time"
)

type Cfg struct {
	Port           string        `mapstructure:"PORT"`
	Host           string        `mapstructure:"HOST"`
	ChromeSetup    string        `mapstructure:"CHROME_SETUP"`
	Leakless       bool          `mapstructure:"LEAKLESS"`
	AnalyzerType   string        `mapstructure:"ANALYZER_TYPE"`
	AnalyzeTimeOut time.Duration `mapstructure:"ANALYZE_TIMEOUT"`
	InMemStoreTTL  time.Duration `mapstructure:"IN_MEM_STORE_TTL"`
	Headless       bool          `mapstructure:"HEADLESS"`
}

var Config *Cfg

func setDefaults() {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("LEAKLESS", true)
	viper.SetDefault("ANALYZER_TYPE", "rod")
	viper.SetDefault("ANALYZE_TIMEOUT", 2)
	viper.SetDefault("IN_MEM_STORE_TTL", 5)
	viper.SetDefault("HEADLESS", true)
}

func GetConfig() *Cfg {
	ctx := context.Background()

	viper.AutomaticEnv()
	bindEnvVariables()

	setDefaults()

	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	_, b, _, _ := runtime.Caller(0)
	viper.AddConfigPath(filepath.Dir(b) + "/..")

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
	_ = viper.BindEnv("ANALYZER_TYPE")
	_ = viper.BindEnv("ANALYZE_TIMEOUT")
	_ = viper.BindEnv("IN_MEM_STORE_TTL")
	_ = viper.BindEnv("HEADLESS")
}
