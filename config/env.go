package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type EnvConfig struct {
	*viper.Viper
}

func NewViperContext(c *viper.Viper) *EnvConfig {
	return &EnvConfig{Viper: c}
}

type ConfigContext interface {
	GetString(key, defaultVal string) string
	GetInt(key string, defaultVal int) int
}

func (c *EnvConfig) GetString(key, defaultVal string) string {
	if key != "" {
		return c.Viper.GetString(key)
	}
	return defaultVal
}

func (c *EnvConfig) GetInt(key string, defaultVal int) int {
	if key != "" {
		return c.Viper.GetInt(key)
	}
	return defaultVal
}

func (c *EnvConfig) InitConfig() {
	switch os.Getenv("ENV") {
	case "dev":
		fmt.Println("dev")
		os.Setenv("ENV", "dev")
		c.Viper.SetConfigName("config_dev")
	case "uat":
		os.Setenv("ENV", "uat")
		c.Viper.SetConfigName("config")
	default:
		c.Viper.SetConfigName("config")
	}

	c.Viper.SetConfigType("yaml")
	c.Viper.AddConfigPath(".")
	c.Viper.AutomaticEnv()
	c.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := c.Viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
