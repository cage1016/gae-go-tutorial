package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gomodule/redigo/redis"

	"github.com/cage1016/gae-lab-001/internal/app/count/endpoints"
	"github.com/cage1016/gae-lab-001/internal/app/count/service"
	"github.com/cage1016/gae-lab-001/internal/app/count/transports"
	pkgredis "github.com/cage1016/gae-lab-001/internal/pkg/redis"
)

const (
	defServiceName string = "count"
	defLogLevel    string = "error"
	defServiceHost string = "localhost"
	defHTTPPort    string = "8180"
	defREDISHost   string = "localhost"
	defREDISPort   string = "6379"

	envServiceName string = "SERVICE_NAME"
	envLogLevel    string = "LOG_LEVEL"
	envServiceHost string = "SERVICE_HOST"
	envHTTPPort    string = "PORT"
	envREDISHost   string = "REDIS_HOST"
	envREDISPort   string = "REDIS_PORT"
)

type config struct {
	serviceName string
	logLevel    string
	serviceHost string
	httpPort    string
	redisConfig pkgredis.Config
}

// Env reads specified environment variable. If no value has been found,
// fallback is returned.
func env(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowInfo())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	cfg := loadConfig(logger)
	logger = log.With(logger, "service", cfg.serviceName)
	level.Info(logger).Log("version", service.Version, "commitHash", service.CommitHash, "buildTimeStamp", service.BuildTimeStamp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := connectToREDIS(cfg.redisConfig, logger).GetContext(ctx)
	if err != nil {
		level.Error(logger).Log(
			"method", "connectToREDIS",
			"host", cfg.redisConfig.Host,
			"port", cfg.redisConfig.Port,
			"err", err,
		)
		os.Exit(1)
	}
	defer conn.Close()

	service := NewServer(conn, logger)
	endpoints := endpoints.New(service, logger)

	wg := &sync.WaitGroup{}

	go startHTTPServer(ctx, wg, endpoints, cfg.httpPort, logger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	cancel()
	wg.Wait()

	fmt.Println("main: all goroutines have told us they've finished")
}

func loadConfig(_ log.Logger) (cfg config) {
	redisConfig := pkgredis.Config{
		Host: env(envREDISHost, defREDISHost),
		Port: env(envREDISPort, defREDISPort),
	}
	cfg.redisConfig = redisConfig
	cfg.serviceName = env(envServiceName, defServiceName)
	cfg.logLevel = env(envLogLevel, defLogLevel)
	cfg.serviceHost = env(envServiceHost, defServiceHost)
	cfg.httpPort = env(envHTTPPort, defHTTPPort)
	return cfg
}

func connectToREDIS(cfg pkgredis.Config, _ log.Logger) *redis.Pool {
	return pkgredis.Connect(cfg)
}

func NewServer(conn redis.Conn, logger log.Logger) service.CountService {
	service := service.New(conn, logger)
	return service
}

func startHTTPServer(ctx context.Context, wg *sync.WaitGroup, endpoints endpoints.Endpoints, port string, logger log.Logger) {
	wg.Add(1)
	defer wg.Done()

	if port == "" {
		level.Error(logger).Log("protocol", "HTTP", "exposed", port, "err", "port is not assigned exist")
		return
	}

	p := fmt.Sprintf(":%s", port)
	// create a server
	srv := &http.Server{Addr: p, Handler: transports.NewHTTPHandler(endpoints, logger)}
	level.Info(logger).Log("protocol", "HTTP", "exposed", port)
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			level.Info(logger).Log("Listen", err)
		}
	}()

	<-ctx.Done()

	// shut down gracefully, but wait no longer than 5 seconds before halting
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ignore error since it will be "Err shutting down server : context canceled"
	srv.Shutdown(shutdownCtx)

	level.Info(logger).Log("protocol", "HTTP", "Shutdown", "http server gracefully stopped")
}
