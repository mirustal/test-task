package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	LogType  string         `yaml:"log"`
	Env      string         `yaml:"env"`
	Postgres PostgresConfig `yaml:"postgres"`
	Rabbitmq   RabbitConfig  `yaml:"rabbitmq`
	REST     REST           `yaml:"rest"`
	MockDB   MockDB         `yaml:"mock_db"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	MaxConn	 int	`yaml:"maxconn"`
}

type RabbitConfig struct {
	Host     string `yaml:host`
	Port     int    `yaml:port`
	User     string `yaml:user`
	Password string `yaml:password`
}

type REST struct {
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type MockDB struct {
	UserCount int `yaml:"user_count"`
	MsgCount  int `yaml:"msg_count"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		panic("config path is empty")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config " + err.Error())
	}

	return &cfg
}
