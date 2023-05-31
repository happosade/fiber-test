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
		AppName:      "CSP Server",
		ServerHeader: "Central bureau of Content Security Policy (CSP) reporting agency",
		Prefork:      true,
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
	})

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
			return c.SendStatus(429)
		},
	}))

	// Monitoring
	prometheus := fiberprometheus.New("CSP_reporting_agency")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	// Path for health check
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})

	// Greetings/testings path
	app.Get("/:name", func(c *fiber.Ctx) error {
		return c.SendString("Hello, " + c.Params("name") + "!")
	})

	// Reporting path
	app.Post("/report", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		c.Response().BodyWriter().Write([]byte("Thanks!"))
		return nil
	})

	app.Listen(":3000")
}
