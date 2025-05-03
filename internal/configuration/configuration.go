package configuration

import (
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Logger    *Logger     `yaml:"Logger"`
	WebServer *HttpServer `yaml:"WebServer"`
	Version   string
}

type Logger struct {
	Level logrus.Level `yaml:"Level"`
}

type HttpServer struct {
	Port              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	ShutdownTimeout   time.Duration
	ReadHeaderTimeout time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
}

func NewConfig() (*Config, error) {

	var conf Config

	viperNew := viper.New()

	viperNew.AddConfigPath(".")
	viperNew.SetConfigName("configuration")
	err := viperNew.ReadInConfig()
	if err != nil {
		return nil, errorspkg.NewErrViperReadInConfig(err)
	}

	err = viperNew.UnmarshalKey("Logger", &conf.Logger)
	if err != nil {
		return nil, errorspkg.NewErrReadConfigViper("Logger", err)
	}

	err = viperNew.UnmarshalKey("HttpServer", &conf.WebServer)
	if err != nil {
		return nil, errorspkg.NewErrReadConfigViper("HttpServer", err)
	}

	return &conf, nil
}
