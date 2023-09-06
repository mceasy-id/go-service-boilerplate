package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"mceasy/service-demo/config"
	"mceasy/service-demo/pkg/apperror"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	goccyjson "github.com/goccy/go-json"
	"github.com/segmentio/encoding/json"
)

type Server struct {
	Config config.Config
	Logger *zap.SugaredLogger
	Fiber  *fiber.App
	DB     *sqlx.DB
}

func NewServer(config config.Config, logger *zap.SugaredLogger, db *sqlx.DB) *Server {
	var fiberConfig fiber.Config
	fiberConfig.ErrorHandler = apperror.HttpHandleError
	fiberConfig.AppName = config.App.Name
	fiberConfig.DisableStartupMessage = true
	fiberConfig.JSONEncoder = goccyjson.Marshal
	fiberConfig.JSONDecoder = json.Unmarshal

	return &Server{
		Config: config,
		Logger: logger,
		Fiber:  fiber.New(fiberConfig),
		DB:     db,
	}
}

func (s *Server) Run() error {
	// Request Logger Middleware
	if s.Config.App.Env == "local" {
		config := logger.ConfigDefault
		config.Format = "[${time}] ${status} ${method} ${path}\n"
		s.Fiber.Use(logger.New(config))
	}

	// Trace Middleware
	if s.Config.Observability.Enable {
		s.Fiber.Use(otelfiber.Middleware())
	}

	// Recover Middleware
	s.Fiber.Use(recover.New(recover.Config{EnableStackTrace: true}))

	// Swagger Handler
	// s.Fiber.Get("/swagger/*", swagger.HandlerDefault)

	// Map App Handlers
	err := s.MapHandlers()
	if err != nil {
		return err
	}

	// Graceful Shutdown
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGINT)
	go func() {
		<-quit
		s.Fiber.Shutdown()
	}()

	// Run Fiber
	s.Logger.Infof("App started")
	return s.Fiber.Listen(fmt.Sprintf(":%s", s.Config.App.Port))
}
