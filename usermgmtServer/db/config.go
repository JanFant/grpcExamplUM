package db

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

var (
	confPath = `usermgmtServer/db/config.toml`
)

func init() {
	var configPath string
	flag.StringVar(&configPath, "config-path", confPath, "path to config file")

	GlobalConfig = NewConfig()

	if _, err := toml.DecodeFile(configPath, &GlobalConfig); err != nil {
		fmt.Println("Can't load config file : ", err.Error())
		os.Exit(1)
	}
}

// GlobalConfig decode config toml
var GlobalConfig *Config

// Config struct all config toml element
type Config struct {
	PSQLConfig ConfigPSQL `toml:"PSQL"`
}

// NewConfig create GlobalConfig
func NewConfig() *Config {
	return &Config{}
}

type ConfigPSQL struct {
	Name         string `toml:"db_name"`
	Pass         string `toml:"db_password"`
	User         string `toml:"db_user"`
	Type         string `toml:"db_type"`
	Host         string `toml:"db_host"`
	Port         string `toml:"db_port"`
	MaxOpenConst int    `toml:"db_SetMaxOpenConst"`
	MaxIdleConst int    `toml:"db_SetMaxIdleConst"`
}

// GetPSQLUrl URL for PSQL connect
func (conf *ConfigPSQL) GetPSQLUrl() string {
	return fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", conf.Host, conf.User, conf.Name, conf.Pass)
}
