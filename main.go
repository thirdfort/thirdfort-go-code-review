package main

import (
	"log/slog"
	"time"

	"github.com/thirdfort/go-slogctx"
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
