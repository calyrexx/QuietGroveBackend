package app

import (
	"context"
	"github.com/calyrexx/zeroslog"
	"github.com/robfig/cron/v3"
	"log/slog"
)

type specFn struct {
	spec string
	fn   cronFunc
}

type cronFunc func(context.Context) error

type AppCron struct {
	cron       *cron.Cron
	fns        []specFn
	fnsOnStart []cronFunc
	logger     *slog.Logger
}

func NewAppCron(logger *slog.Logger) (*AppCron, error) {
	l := logger.With(zeroslog.ServiceKey, "AppCron")

	m := make([]specFn, 0, 16)
	mOnStart := make([]cronFunc, 0, 16)

	cr := cron.New(cron.WithSeconds())
	return &AppCron{
		cron:       cr,
		fns:        m,
		logger:     l,
		fnsOnStart: mOnStart,
	}, nil
}

func (c *AppCron) Add(spec []string, fn cronFunc) {
	for _, sp := range spec {
		c.fns = append(c.fns, specFn{sp, fn})
	}
}

func (c *AppCron) AddOnStart(fn cronFunc) {
	c.fnsOnStart = append(c.fnsOnStart, fn)
}

func (c *AppCron) Start(ctx context.Context) error {
	for _, v := range c.fns {
		spec, fn := v.spec, v.fn
		f, err := c.cron.AddFunc(spec, func() {
			err := fn(ctx)
			if err != nil {
				c.logger.Error(err.Error())
			}
		})
		if err != nil {
			return err
		}
		c.cron.Entry(f)
	}

	c.logger.Info("AppCron has been started!")

	for _, fn := range c.fnsOnStart {
		err := fn(ctx)
		if err != nil {
			c.logger.Error(err.Error())
		}
	}

	c.cron.Start()
	defer c.cron.Stop()

	<-ctx.Done()
	return nil
}
