package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// AppConfig a bag containing all necessary things for this app.
type AppConfig struct {
	Config *viper.Viper
	InfL   *log.Logger
	ErrL   *log.Logger
}

// InitConfig init config and return preconfigured viper instance.
func InitConfig(filePath string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(filePath)
	v.SetConfigName("app")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

// SetupDefault setup default value and return error if required fields is not present.
func SetupDefault(v *viper.Viper) error {
	// global
	if !v.IsSet("root") {
		return fmt.Errorf("`root` for this app root directories is required")
	}
	v.SetDefault("max_days", 6)

	// databases
	v.SetDefault("db.max_worker", 1)
	if !v.IsSet("db.max_days") {
		v.SetDefault("db.max_days", v.GetInt("max_days"))
	}

	// apps/directories
	v.SetDefault("app.max_worker", 1)
	if !v.IsSet("app.max_days") {
		v.SetDefault("app.max_days", v.GetInt("max_days"))
	}

	return nil
}
