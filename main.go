package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/thirdfort/go-slogctx"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/thirdfort/thirdfort-go-code-review/internal/cache"
	"github.com/thirdfort/thirdfort-go-code-review/internal/config"
	"github.com/thirdfort/thirdfort-go-code-review/internal/repositories"
	"github.com/thirdfort/thirdfort-go-code-review/internal/service"
	"github.com/thirdfort/thirdfort-go-code-review/internal/service/web"
)

func main() {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}

	logger := initLogger(conf)

	ds, err := repositories.NewStore(conf, logger)
	if err != nil {
		panic(err)
	}
	defer ds.Close()

	cache := cache.New()

	service, err := service.New(
		conf,
		logger,
		ds,
		cache,
	)
	if err != nil {
		panic(err)
	}

	web, err := web.New(conf, service)
	if err != nil {
		panic(err)
	}

	sentry.CaptureException(fmt.Errorf("Consumer API started - version: %s - build date: %s - git hash: %s", internal.Version, internal.BuildTime, internal.GitHash))

	err = web.Run()
	if err != nil {
		panic(err)
	}
}

func initLogger(conf *config.Config) *slogctx.Logger {
	l := slogctx.New(&slogctx.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelInfo,
		ReplaceAttr: slogctx.ReplaceAttr,
		TimeFormat:  "",
	})

	if conf.App.Debug {
		l = slogctx.New(&slogctx.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelDebug,
			ReplaceAttr: slogctx.ReplaceAttr,
			TimeFormat:  time.DateTime,
		})
	}

	return l
}
