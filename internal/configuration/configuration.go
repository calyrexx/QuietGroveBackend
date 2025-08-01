package configuration

import (
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"log/slog"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Logger       *Logger       `yaml:"Logger"`
		WebServer    *HttpServer   `yaml:"WebServer"`
		AppCron      *AppCron      `yaml:"AppCron"`
		Reservations *Reservations `yaml:"Reservations"`
		Version      string
	}

	AppCron struct {
		UpdateReservationsStatuses CronConfig
		GetForReminder             CronConfig
	}

	CronConfig struct {
		Spec []string
	}

	Logger struct {
		Level slog.Level `yaml:"Level"`
	}

	Reservations struct {
		PriceCoefficients     []PriceCoefficient
		NotificationThreshold int
	}

	PriceCoefficient struct {
		Start time.Time
		End   time.Time
		Rate  float64
	}

	PriceCoefficientTemp struct {
		Start string
		End   string
		Rate  float64
	}

	HttpServer struct {
		Port              string
		ReadTimeout       time.Duration
		WriteTimeout      time.Duration
		ShutdownTimeout   time.Duration
		ReadHeaderTimeout time.Duration
		IdleTimeout       time.Duration
		MaxHeaderBytes    int
	}
)

func NewConfig() (*Config, error) {
	var (
		conf Config
		temp struct {
			NotificationThreshold int
			PriceCoefficients     []PriceCoefficientTemp
		}
	)

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

	err = viperNew.UnmarshalKey("AppCron", &conf.AppCron)
	if err != nil {
		return nil, errorspkg.NewErrReadConfigViper("AppCron", err)
	}

	err = viperNew.UnmarshalKey("Reservations", &temp)
	if err != nil {
		return nil, errorspkg.NewErrReadConfigViper("PriceCoefficients", err)
	}

	pc, err := parseDate(temp.PriceCoefficients)
	if err != nil {
		return nil, errorspkg.NewErrReadConfigViper("PriceCoefficients", err)
	}

	conf.Reservations = &Reservations{
		NotificationThreshold: temp.NotificationThreshold,
		PriceCoefficients:     pc,
	}

	return &conf, nil
}

func parseDate(req []PriceCoefficientTemp) ([]PriceCoefficient, error) {
	result := make([]PriceCoefficient, 0, len(req))
	for _, pc := range req {
		start, pErr := time.Parse("2006-01-02", pc.Start)
		if pErr != nil {
			return nil, pErr
		}
		end, pErr := time.Parse("2006-01-02", pc.End)
		if pErr != nil {
			return nil, pErr
		}
		result = append(result, PriceCoefficient{
			Start: start,
			End:   end,
			Rate:  pc.Rate,
		})
	}
	return result, nil
}
