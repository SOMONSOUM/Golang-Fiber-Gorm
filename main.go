package main

import (
	"log"
	"os"

	"github.com/SOMONSOUM/go-fiber/config"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var db *gorm.DB = config.SetupDatabaseConnection()

func main() {
	defer config.CloseDatabaseConnection(db)
	app := fiber.New()

	authRoutes := app.Group("api/auth")
	{
		authRoutes.Get("/users", func(c *fiber.Ctx) error {
			return c.JSON("All users")
		})
		authRoutes.Get("/register", func(c *fiber.Ctx) error {
			return c.JSON("Signup")
		})
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello, World👋!")
	})

	log.Fatal(app.Listen(os.Getenv("PORT")))
}
