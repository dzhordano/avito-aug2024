package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var (
	cfgPath = flag.String("c", "", "config path")
)

type Config struct {
	Env      string         `yaml:"env"`
	Postgres PostgresConfig `yaml:"postgres"`
	HTTP     HTTPConfig     `yaml:"http"`
	Auth     AuthConfig     `yaml:"auth"`
}

type PostgresConfig struct {
	Host     string `env:"PGHOST"`
	Port     int    `env:"PGPORT"`
	User     string `env:"PGUSER"`
	Password string `env:"PGPASS"`
	Database string `env:"PGDB"`
	SSLMode  string `env:"PGSSLMODE"`
}

type HTTPConfig struct {
	Host               string        `yaml:"host"`
	Port               string        `yaml:"port"`
	ReadTimeout        time.Duration `yaml:"rd_timeout"`
	WriteTimeout       time.Duration `yaml:"wr_timeout"`
	MaxHeaderMegabytes int           `yaml:"max_header_megabytes"`
}

type AuthConfig struct {
	SecretKey string        `env:"AUTH_SECRET_KEY"`
	TokenTTL  time.Duration `yaml:"token_ttl"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func MustLoad() *Config {
	flag.Parse()
	path := *cfgPath

	if path == "" {
		path = os.Getenv("DEFAULT_CONFIG_PATH")
	}

	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}
