package main

import (
	"os"
	"strconv"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/happosade/fiber-test/elastic"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func init() {
	elastic.ConnectES()
}

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "Reporting Server",
		ServerHeader: "Central bureau of Browser reported metrics",
		Prefork:      true,
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
	})	
	
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	// Ratelimit
	i, _ := strconv.Atoi(getEnv("RATE_LIMIT", "30"))
	app.Use(limiter.New(limiter.Config{
		Max:               i,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("X-Forwarded-For")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	}))

	// Monitoring
	prometheus := fiberprometheus.New("reporting_server")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)


	// Path for health check
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})

	app.Get("/healthz", func(c *fiber.Ctx) error {
		if elastic.ConnectStatus() {
			return c.SendStatus(fiber.StatusOK)
		}

		return c.SendStatus(fiber.StatusServiceUnavailable)
	})

	// Greetings/testings path
	app.Get("/:name", func(c *fiber.Ctx) error {
		return c.SendString("Hello, " + c.Params("name") + "!")
	})

	// Reporting CSP
	app.Post("/csp", func(c *fiber.Ctx) error {
		c.Accepts(fiber.MIMEApplicationJSONCharsetUTF8)
		c.Response().BodyWriter().Write([]byte("Thanks!"))
		return nil
	})

	app.Listen(":3000")
}
