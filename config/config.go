package env

import (
	"crypto/rsa"
)

var settingCfg *Config

var (
	AppName     string
	Version     string
	BuildTime   string
	BuildCommit string
)

type Config struct {

	//ENV
	Env string

	// App
	AppName     string
	Version     string
	BuildTime   string
	BuildCommit string
	Production  bool
	Debug       bool

	// Resource
	// Db           database.Database
	DbHost     string
	DbPort     string
	DbUsername string
	DbPass     string
	DbName     string
	// Cache        cache.Redis
	CacheHost    string
	CachPort     string
	CachStore    int
	ServerConfig ServerConfig

	//server
	PrivateKey *rsa.PrivateKey
}

type ServerConfig struct {
	Host string
	Port string
}

func NewCfg() *Config {
	return &Config{
		AppName:     AppName,
		Version:     Version,
		BuildTime:   BuildTime,
		BuildCommit: BuildCommit,
	}
}

func GetCfg() *Config {
	if settingCfg == nil {
		settingCfg = NewCfg()
	}
	return settingCfg
}

func (cfg *Config) Load(env *EnvConfig) error {
	cfg.ServerConfig = ServerConfig{
		Host: env.GetString("app.url", "127.0.0.1"),
		Port: env.GetString("app.port", "7002"),
	}

	cfg.DbHost = env.GetString("pg.host", "127.0.0.1")
	cfg.DbPort = env.GetString("pg.port", "127.0.0.1")
	cfg.DbUsername = env.GetString("pg.username", "127.0.0.1")
	cfg.DbPass = env.GetString("pg.password", "127.0.0.1")
	cfg.DbName = env.GetString("pg.database", "127.0.0.1")

	cfg.CachPort = env.GetString("cache.port", "")
	cfg.CacheHost = env.GetString("cache.host", "")
	cfg.CachStore = env.GetInt("cache.db", 0)

	cfg.Env = env.GetString("env", "dev")
	return nil
}
