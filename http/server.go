package http

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"order-backend/config"
	"order-backend/http/handler"
	"order-backend/http/router"
	"order-backend/order"
	"order-backend/pkg/logger"
	"order-backend/pkg/postgres"
	"order-backend/repository"
	"os"
	"os/signal"
	"syscall"
)

func StartServer() {
	log := logger.Log()
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to get config. got error %s", err.Error()))
	}

	db, err := postgres.ConnectDB(config.Postgres)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to init database got error %s", err.Error()))
	}

	if config.Postgres.MigrationOnStartupEnabled {
		err = postgres.MigrateUp(db, config.Postgres.MigrationFilePath)
		if err != nil {
			log.Fatal(fmt.Sprintf("failed to migrate database with error %s", err.Error()))
		}
	}

	orderRepository := repository.NewOrderRepository(db)

	orderService := order.NewOrderService(orderRepository)

	orderHandler := handler.NewOrderHandler(orderService)

	appServer := initFiberApp()
	router.SetupRoutes(appServer, router.Opts{OrderHandler: orderHandler})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		doneChan := make(chan bool, 1)

		<-sigint
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), config.App.GracefulShutdownTimeout)
		defer shutdownCancel()

		go func() {
			log.Info("Shutdown signal received")

			// Get postgres database connection
			if err := db.Close(); err != nil {
				log.Error(fmt.Sprintf("Failed shutdown Postgres database. Error: %v", err))
			} else {
				log.Info("Successfully shutdown Postgres database")
			}

			// Shutdown http server
			if err := appServer.Shutdown(); err != nil {
				log.Error(fmt.Sprintf("Failed to shutdown HTTP server. Error: %v", err))
			} else {
				log.Info("Successfully shutdown HTTP server")
			}

			doneChan <- true
		}()

		select {
		case <-shutdownCtx.Done():
			log.Info("Graceful shutdown timed out. Some tasks may not have completed.")
			return
		case <-doneChan:
			log.Info("All connections have been shut down gracefully. Exiting...")
			return
		}
	}()

	if err := appServer.Listen(config.App.Port); err != nil {
		log.Error(fmt.Sprintf("Error starting HTTP server: %v", err))
	}

}

func initFiberApp() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	return app
}
