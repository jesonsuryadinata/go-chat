package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/pusher/pusher-http-go/v5"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	pusherClient := pusher.Client{
		AppID:   "1648273",
		Key:     "c6998e24f5996790c71d",
		Secret:  "77a1884f6afc83c9dc72",
		Cluster: "ap1",
		Secure:  true,
	}

	app.Post("/api/messages", func(c *fiber.Ctx) error {
		var data map[string]string

		if err := c.BodyParser(&data); err != nil {
			return err
		}

		pusherClient.Trigger("chat", "message", data)

		return c.JSON([]string{})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world !")
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}
