// @title           fiber-skeleton API
// @version         1.0
// @description     A skeleton API built with Go Fiber.
// @host            localhost:3000
// @BasePath        /

package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	fiberrequestid "github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/jhseoeo/fiber-skeleton/src/config"
	_ "github.com/jhseoeo/fiber-skeleton/docs"
	"github.com/jhseoeo/fiber-skeleton/src/handler"
	"github.com/jhseoeo/fiber-skeleton/src/middleware"
	pkglog "github.com/jhseoeo/fiber-skeleton/src/pkg/log"
	"github.com/jhseoeo/fiber-skeleton/src/repository"
	"github.com/jhseoeo/fiber-skeleton/src/service"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/swag"
)

const shutdownTimeout = 10 * time.Second

func main() {
	cfg := config.Load()
	pkglog.Init(cfg)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.NewErrorHandler(),
		BodyLimit:    4 * 1024 * 1024, // 4 MB
	})

	app.Use(middleware.NewRecoverer())
	app.Use(middleware.NewMetrics(app))
	app.Use(buildCORS(cfg))
	app.Use(middleware.NewSecurity())
	app.Use(middleware.NewTimeout(cfg.RequestTimeout))
	app.Use(middleware.NewLogger())
	app.Use(fiberrequestid.New())

	// pprof: only expose in development
	if cfg.Env == "development" {
		app.Use(pprof.New())
		logrus.Info("pprof enabled at /debug/pprof")
	}

	app.Get("/health/live", handler.Liveness)
	app.Get("/health/ready", handler.Readiness)

	// Swagger: GET /swagger/doc.json + minimal UI
	registerSwagger(app)

	exampleRepository := repository.NewExampleRepository()
	exampleService := service.NewExampleService(exampleRepository)
	exampleHandler := handler.NewExampleHandler(exampleService)
	exampleHandler.RegisterRoutes(app)

	// TODO: protect routes with JWT:
	//   api := app.Group("/api", middleware.NewAuthMiddleware([]byte(cfg.JWTSecret)))
	//   exampleHandler.RegisterRoutes(api)

	// TODO: apply rate limiting to specific route groups, e.g.:
	//   import "github.com/gofiber/fiber/v3/middleware/limiter"
	//   api := app.Group("/api", middleware.NewRateLimiter(limiter.Config{Max: 100, Expiration: time.Minute}))
	//   exampleHandler.RegisterRoutes(api)

	// Catch-all: return a JSON 404 for unregistered routes.
	app.Use(handler.NotFound)

	quit := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		logrus.Info("shutting down server...")
		done := make(chan error, 1)
		go func() { done <- app.Shutdown() }()
		select {
		case err := <-done:
			if err != nil {
				logrus.WithError(err).Error("failed to shutdown server")
			}
		case <-time.After(shutdownTimeout):
			logrus.Error("server shutdown timed out")
		}
		close(quit)
	}()

	logrus.Infof("server starting on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		logrus.WithError(err).Fatal("failed to start server")
	}
	<-quit
}

func buildCORS(cfg *config.Config) fiber.Handler {
	if cfg.CORSAllowOrigins == "" {
		if cfg.Env == "production" {
			logrus.Warn("CORS_ALLOW_ORIGINS is not set; allowing all origins in production is insecure")
		}
		return middleware.NewCORS()
	}
	return middleware.NewCORS(cors.Config{AllowOrigins: []string{cfg.CORSAllowOrigins}})
}

func registerSwagger(app *fiber.App) {
	// Serve swagger JSON spec.
	// Run `swag init -g cmd/main.go` to regenerate docs from annotations.
	app.Get("/swagger/doc.json", func(c fiber.Ctx) error {
		doc, err := swag.ReadDoc()
		if err != nil {
			return err
		}
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return c.SendString(doc)
	})

	// Serve a minimal Swagger UI (loads swagger-ui from CDN).
	app.Get("/swagger", func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return c.SendString(swaggerUIHTML)
	})
}

const swaggerUIHTML = `<!DOCTYPE html>
<html>
<head>
  <title>fiber-skeleton API</title>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>
window.onload = function() {
  SwaggerUIBundle({ url: "/swagger/doc.json", dom_id: '#swagger-ui' });
}
</script>
</body>
</html>`
