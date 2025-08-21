package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ebadidev/arch-node/internal/config"
	"github.com/ebadidev/arch-node/internal/database"
	"github.com/ebadidev/arch-node/internal/http/handlers"
	v1 "github.com/ebadidev/arch-node/internal/http/handlers/v1"
	"github.com/ebadidev/arch-node/pkg/http/middleware"
	"github.com/ebadidev/arch-node/pkg/http/validator"
	"github.com/ebadidev/arch-node/pkg/logger"
	"github.com/ebadidev/arch-node/pkg/xray"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	engine   *echo.Echo
	config   *config.Config
	xray     *xray.Xray
	database *database.Database
	l        *logger.Logger
}

// Run defines the required HTTP routes and starts the HTTP Server.
func (s *Server) Run() {
	s.engine.Use(echoMiddleware.CORS())
	s.engine.Use(middleware.Logger(s.l))
	s.engine.Use(middleware.General())

	s.engine.GET("/", handlers.HomeShow())

	g2 := s.engine.Group("/v1")
	g2.Use(middleware.Authorize(func() string {
		return s.database.Data.Settings.HttpToken
	}))

	g2.GET("/stats", v1.StatsShow(s.xray))
	g2.POST("/configs", v1.ConfigsStore(s.xray))
	g2.POST("/manager", v1.ManagerStore(s.database))

	go func() {
		address := fmt.Sprintf("%s:%d", "0.0.0.0", s.database.Data.Settings.HttpPort)
		if err := s.engine.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.l.Fatal("http server: cannot start", zap.String("address", address), zap.Error(err))
		}
	}()
}

// Close closes the HTTP Server.
func (s *Server) Close() error {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.engine.Shutdown(c); err != nil {
		return errors.WithStack(err)
	}

	s.l.Debug("http server: closed successfully")
	return nil
}

// New creates a new instance of HTTP Server.
func New(config *config.Config, l *logger.Logger, x *xray.Xray, d *database.Database) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Validator = validator.New()

	return &Server{engine: e, config: config, l: l, xray: x, database: d}
}
