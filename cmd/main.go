package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/jhseoeo/fiber-skeleton/src/handler"
	"github.com/jhseoeo/fiber-skeleton/src/middleware"
	"github.com/jhseoeo/fiber-skeleton/src/repository"
	"github.com/jhseoeo/fiber-skeleton/src/service"
	"github.com/sirupsen/logrus"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.NewErrorHandler(),
	})

	app.Use(middleware.NewRecoverer())
	app.Use(middleware.NewLogger())

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	exampleRepository := repository.NewExampleRepository()
	exampleService := service.NewExampleService(exampleRepository)
	exampleHandler := handler.NewExampleHandler(exampleService)
	exampleHandler.RegisterRoutes(app)

	quit := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		logrus.Info("shutting down server...")
		if err := app.Shutdown(); err != nil {
			logrus.WithError(err).Error("failed to shutdown server")
		}
		close(quit)
	}()

	if err := app.Listen(":3000"); err != nil {
		logrus.WithError(err).Fatal("failed to start server")
	}
	<-quit
}
