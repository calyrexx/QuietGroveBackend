package configuration

import (
	"github.com/spf13/viper"
	"time"
)

type Credentials struct {
	Postgres    Postgres
	TelegramBot TelegramBot
}

type Postgres struct {
	Host               string        `yaml:"Host"`
	Port               string        `yaml:"Port"`
	Database           string        `yaml:"Database"`
	User               string        `yaml:"User"`
	Password           string        `yaml:"Password"`
	MinConnections     int32         `yaml:"MinConnections"`
	MaxConnections     int32         `yaml:"MaxConnections"`
	IdleConnection     time.Duration `yaml:"IdleConnection"`
	LifeTimeConnection time.Duration `yaml:"LifeTimeConnection"`
	JitterConnection   time.Duration `yaml:"JitterConnection"`
}

type TelegramBot struct {
	Token        string `yaml:"Token"`
	AdminChatIDs []int  `yaml:"AdminChatID"`
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

	err = viperNew.UnmarshalKey("TelegramBot", &creds.TelegramBot)
	if err != nil {
		return nil, err
	}

	return &creds, nil
}
