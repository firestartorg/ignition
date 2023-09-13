package maintenance

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"gitlab.com/firestart/ignition"
	"net/http"
	"os"
	"time"
)

type ServerConfig struct {
	Metrics            bool
	Health             bool
	HealthCheckTimeout int
}

type Server struct {
	config ServerConfig
	logger zerolog.Logger

	// health checks are used to determine if the server is healthy and ready
	healthLive  *HealthCheck
	healthReady *HealthCheck

	router *httprouter.Router
}

func NewServer(config *ServerConfig, logger *zerolog.Logger) *Server {
	if config == nil {
		config = &ServerConfig{
			Metrics: true,
			Health:  true,
		}
	}

	if logger == nil {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger = ignition.Reference(zerolog.New(os.Stderr).With().Timestamp().Logger())
	}

	srv := &Server{
		config: *config,
		logger: *logger,
	}

	// initialize router and add health and metrics endpoints
	router := httprouter.New()
	if config.Health {
		router.GET("/health/ready", srv.readinessProbeHandle)
		router.GET("/health/live", srv.livenessProbeHandle)
	}
	if config.Metrics {
		router.Handler(http.MethodGet, "/metrics", promhttp.Handler())
	}
	srv.router = router

	return srv
}

func (s *Server) SetLivenessProbe(healthCheck HealthCheck) {
	s.healthLive = &healthCheck
}

func (s *Server) SetReadinessProbe(healthCheck HealthCheck) {
	s.healthReady = &healthCheck
}

func (s *Server) probeHandle(check *HealthCheck, writer http.ResponseWriter, request *http.Request) {
	// if no health check is set, return ok
	if check == nil {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("Ok"))
		return
	}

	// create a context with a timeout for the health check
	ctx, cancel := context.WithTimeout(request.Context(), time.Duration(s.config.HealthCheckTimeout)*time.Second)
	defer cancel()

	// run the health check
	ok, err := (*check)(ctx)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Health check failed")

		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("Failed"))
		return
	}

	if !ok {
		writer.WriteHeader(http.StatusServiceUnavailable)
		_, _ = writer.Write([]byte("Unavailable"))
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (s *Server) readinessProbeHandle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	s.probeHandle(s.healthReady, writer, request)
}

func (s *Server) livenessProbeHandle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	s.probeHandle(s.healthLive, writer, request)
}

var _ http.Handler = (*Server)(nil)

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) ListenAndServe(addr string) {
	s.logger.Info().
		Str("addr", addr).
		Msg("Maintenance server listening")

	err := http.ListenAndServe(addr, s)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Maintenance server listener error")

		return
	}
}
