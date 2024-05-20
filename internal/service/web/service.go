package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	sloggin "github.com/samber/slog-gin"
	"github.com/thirdfort/go-gin/middleware"
	"github.com/thirdfort/go-otel/instrumentation/otelgin"
	"github.com/thirdfort/go-slogctx"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/thirdfort/thirdfort-go-code-review/internal/config"
	"github.com/thirdfort/thirdfort-go-code-review/internal/service"
	"github.com/wI2L/fizz"
)

type WebService struct {
	Config *config.Config
	Fizz   *fizz.Fizz
	Gin    *gin.Engine
	srv    *service.Service
	server *http.Server
}

func New(conf *config.Config, srv *service.Service) (*WebService, error) {
	if !conf.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	logConfig := sloggin.Config{
		WithUserAgent:      false,
		WithRequestID:      false,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,
	}

	engine := gin.New()

	if conf.Otel.Enabled {
		engine.Use(
			otelgin.Metrics(),
			otelgin.Trace(fmt.Sprintf("%s-otel", conf.App.Name)),
		)
	}

	engine.Use(
		sloggin.NewWithConfig(srv.Logger.GetLogger(), logConfig),
		gin.Recovery(),
		middleware.Sentry(),
	)

	engine.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error": gin.H{
				"message": "Not found",
			},
		})
	})

	tonic.SetErrorHook(ErrHook)

	f := fizz.NewFromEngine(engine)

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", conf.Service.Address, conf.Service.Port),
		ReadHeaderTimeout: 15 * time.Second,
		Handler:           f,
	}

	return &WebService{
		Config: conf,
		Fizz:   f,
		Gin:    engine,
		srv:    srv,
		server: server,
	}, nil
}

func (s *WebService) Run() error {
	slogctx.Info(context.TODO(), "Consumer API started!",
		slog.String("version", internal.Version),
		slog.String("version", internal.BuildTime),
		slog.String("version", internal.GitHash),
	)

	s.setV1Routes()

	// Non-public routes
	s.setHealth()
	s.setOpenAPI()

	return s.server.ListenAndServe()
}
