package httpserver

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router interface {
	Attach(*echo.Group)
}

type Options struct {
	Addr             string
	MaxBodySize      string
	CompressionLevel int
	Logger           *slog.Logger
}

type Server struct {
	echo *echo.Echo
	addr string
}

func New(opts *Options) *Server {
	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogLatency:   true,
		LogMethod:    true,
		LogError:     true,
		LogStatus:    true,
		LogRequestID: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			msg := v.Method + " " + v.URI
			if v.Error != nil {
				opts.Logger.Error(msg, "err", v.Error, "status", v.Status, "requestID", v.RequestID, "latency", v.Latency)
				return nil
			}

			opts.Logger.Info(msg, "status", v.Status, "requestID", v.RequestID, "latency", v.Latency)
			return nil
		},
	}))
	e.Use(middleware.BodyLimit(opts.MaxBodySize))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: opts.CompressionLevel}))
	e.Use(middleware.CORS())

	return &Server{
		echo: e,
		addr: opts.Addr,
	}
}

func (s *Server) AddRouter(prefix string, r Router) {
	g := s.echo.Group(prefix)
	r.Attach(g)
}

func (s *Server) Start() error {
	return s.echo.Start(s.addr)
}
