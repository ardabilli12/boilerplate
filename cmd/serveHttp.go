package cmd

import (
	"boilerplate-service/config"
	"boilerplate-service/pkg/logger"
	"boilerplate-service/pkg/mySqlExt"
	"boilerplate-service/pkg/newRelicExt"
	"boilerplate-service/pkg/redisExt"
	"boilerplate-service/port/http"
	"context"
	"fmt"
	"log"
	netHttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	healthCheckRepo "boilerplate-service/internal/repository/healthCheck"
	healthCheckSvc "boilerplate-service/internal/service/v1/healthCheck"
	v1HealthCheckController "boilerplate-service/port/http/controller/v1/healthCheck"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveHttpCmd)
}

var serveHttpCmd = &cobra.Command{
	Use:   "serveHttp",
	Short: "Start HTTP server",
	Long:  `Start Boilerplate HTTP server`,
	Run: func(cmd *cobra.Command, args []string) {
		// Init config
		config, secret, err := config.LoadConfig(cfgFile, scrtFile)
		if err != nil {
			log.Fatalf("Unable to load configuration and secret: %v", err)
		}

		// Logger
		loggerConfig := logger.Config{
			Environment: config.Environment,
			ServiceName: config.ServiceName,
		}
		logger, err := logger.New(
			loggerConfig,
		)
		if err != nil {
			fmt.Printf("Unable to init logger, %v", err)
			panic(err)
		}
		defer logger.Sync()

		// New Relic
		newRelicExtConfig := newRelicExt.Config{
			Environment: config.Environment,
			LicenseKey:  secret.NewRelicLicenseKey,
			ServiceName: config.ServiceName,
			Logger:      logger,
		}
		newRelic, err := newRelicExt.New(newRelicExtConfig)
		if err != nil {
			fmt.Printf("Unable to init new relic, %v", err)
			panic(err)
		}

		defer newRelic.Shutdown(10 * time.Second)

		// MySQL Database
		mysqlExtConfig := mySqlExt.Config{
			Host:         config.MySQLConfig.Host,
			Port:         config.MySQLConfig.Port,
			Username:     secret.MySQLSecret.Username,
			Password:     secret.MySQLSecret.Password,
			DBName:       secret.MySQLSecret.Database,
			MaxIdleConns: config.MySQLConfig.MaxIdleConns,
			MaxIdleTime:  config.MySQLConfig.MaxOpenConns,
			MaxLifeTime:  config.MySQLConfig.MaxLifeTime,
			MaxOpenConns: config.MySQLConfig.MaxOpenConns,
		}

		dbClient, err := mySqlExt.New(mysqlExtConfig)
		if err != nil {
			fmt.Printf("Unable to init mysql gateway, %v", err)
			panic(err)
		}
		defer dbClient.Close()

		// Redis
		cacheClient, err := redisExt.New(
			config.RedisConfig.Host,
			config.RedisConfig.Port,
			secret.RedisSecret.Password,
			config.RedisConfig.CacheDB,
		)
		if err != nil {
			fmt.Printf("Unable to init redis cache, %v", err)
			panic(err)
		}
		defer cacheClient.Close()

		// todo Validator will used later when standardize validation request & response message done
		// validate := validatorExt.New()

		// Init repository
		// e.g. database, external/internal services repository, etc.
		healthCheckRepository := healthCheckRepo.New(
			logger,
			dbClient,
			cacheClient,
		)

		// Init services
		healthCheckService := healthCheckSvc.New(
			config,
			healthCheckRepository,
		)

		// Init controller
		healthCheckController := v1HealthCheckController.New(
			healthCheckService,
		)

		// Init router
		r := http.HttpRoute(
			newRelic,
			logger,
			healthCheckController,
		)

		server := &netHttp.Server{Addr: ":3000", Handler: r}
		// Start the server in a goroutine
		go func() {
			if err := server.ListenAndServe(); err != nil && err != netHttp.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		// Set up signal capturing
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Block until we receive our signal.
		<-quit

		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline or until all connections have returned.
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server Shutdown Failed:%+v", err)
		}

		log.Println("Server gracefully stopped")
	},
}
