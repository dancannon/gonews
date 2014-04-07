package config

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

type Env string

const (
	Dev  string = "development"
	Prod string = "production"
	Test string = "testing"
)

type Config struct {
	Name    string `toml:"name"`
	Env     string `toml:"env"`
	Address string `toml:"address"`

	PublicPath   string `toml:"public_path"`
	TemplatePath string `toml:"template_path"`

	Security  Security          `toml:"security"`
	RethinkDB RethinkDB         `toml:"rethinkdb"`
	Redis     Redis             `toml:"redis"`
	NSQ       NSQ               `toml:"nsq"`
	Workers   map[string]Worker `toml:"workers"`
	GoFetch   GoFetch           `toml:"gofetch"`
}

func (conf *Config) Validate() error {
	if err := conf.NSQ.Validate(); err != nil {
		return err
	}

	return nil
}

type Security struct {
	Secret   string `toml:"secret"`
	BlockKey string `toml:"block_key"`
}

type RethinkDB struct {
	MaxIdle  int    `toml:"max_idle"`
	Address  string `toml:"address"`
	Database string `toml:"database"`
}

type Redis struct {
	MaxIdle     int           `toml:"max_idle"`
	IdleTimeout time.Duration `toml:"idle_timeout"`
	Network     string        `toml:"network"`
	Address     string        `toml:"address"`
	Password    string        `toml:"password"`
}

type NSQ struct {
	MaxInFlight     int      `toml:"max_in_flight"`
	NSQDAddresses   []string `toml:"nsqd_addresses"`
	LookupAddresses []string `toml:"lookup_addresses"`
}

func (conf *NSQ) Validate() error {
	if len(conf.NSQDAddresses) == 0 && len(conf.LookupAddresses) == 0 {
		return fmt.Errorf("nsqd_addresses or lookup_addresses required")
	}

	if len(conf.NSQDAddresses) > 0 && len(conf.LookupAddresses) > 0 {
		return fmt.Errorf("Use only nsqd_addresses or lookup_addresses not both")
	}

	return nil
}

type Worker struct {
	Topic   string
	Channel string
}

type GoFetch struct {
	ConfigFile string `toml:"config_file"`
}

func NewConfig() *Config {
	return &Config{
		Env:     Dev,
		Address: ":3000",
		RethinkDB: RethinkDB{
			MaxIdle: 10,
		},
		Redis: Redis{
			Network:     "tcp",
			MaxIdle:     10,
			IdleTimeout: 240 * time.Second,
		},
		NSQ: NSQ{
			MaxInFlight: 200,
		},
	}
}

func LoadFile(conf *Config, file string) error {
	if _, err := toml.DecodeFile(file, conf); err != nil {
		return err
	}

	if err := conf.Validate(); err != nil {
		return err
	}

	return nil
}
