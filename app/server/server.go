package server

import (
	"creditlimit-connector/app/configs"
	"creditlimit-connector/app/log"
	"creditlimit-connector/app/errors"
	"creditlimit-connector/app/middlewares"
	"creditlimit-connector/app/routes"
	"fmt"
	"os"

	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Serve() *fiber.App {

	app := fiber.New(fiber.Config{
		ErrorHandler: errors.CustomErrorHandler,
	})

	// Middlewares
	app.Use(healthcheck.New())
	app.Use(recover.New())
	app.Use(middlewares.RequestContextMiddleware())

	// Routes
	routes.InitRoutes(app)

	// Start the server in a goroutine preventing it from blocking the main thread
	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", configs.Conf.App.Port)); err != nil {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()

	// Create a channel to listen for OS signals (e.g., SIGINT, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for an interrupt signal (Ctrl+C, SIGTERM, etc.)
	sig := <-sigChan
	log.Infof("Received signal: %s", sig)

	// Initiate graceful shutdown of the HTTP server
	log.Info("Shutting down the server...")
	if err := app.Shutdown(); err != nil {
		log.Errorf("Error during graceful shutdown:", err.Error())
	} else {
		log.Info("Server gracefully shut down")
	}

	return app
}
