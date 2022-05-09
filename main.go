package main

import (
	"log"

	"github.com/SOMONSOUM/go-fiber/config"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var db *gorm.DB = config.SetupDatabaseConnection()

func main() {
	defer config.CloseDatabaseConnection(db)
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello, World👋!")
	})

	log.Fatal(app.Listen(":3000"))
}
