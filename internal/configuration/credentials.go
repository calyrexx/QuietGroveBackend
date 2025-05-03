package configuration

import (
	"github.com/spf13/viper"
	"time"
)

type Credentials struct {
	Postgres Postgres
}

type Postgres struct {
	Host               string        `yaml:"Host"`
	Port               string        `yaml:"Port"`
	Database           string        `yaml:"Database"`
	User               string        `yaml:"User"`
	Password           string        `yaml:"Password"`
	MinConnections     int           `yaml:"MinConnections"`
	MaxConnections     int           `yaml:"MaxConnections"`
	IdleConnection     time.Duration `yaml:"IdleConnection"`
	LifeTimeConnection time.Duration `yaml:"LifeTimeConnection"`
	JitterConnection   time.Duration `yaml:"JitterConnection"`
}

func NewCredentials() (*Credentials, error) {
	var creds Credentials

	viperNew := viper.New()

	viperNew.AddConfigPath(".")
	viperNew.SetConfigName("credentials")

	err := viperNew.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viperNew.UnmarshalKey("Postgres", &creds.Postgres)
	if err != nil {
		return nil, err
	}

	return &creds, nil
}
