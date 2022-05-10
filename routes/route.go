package routes

import (
	"github.com/SOMONSOUM/go-fiber/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello, World👋!")
	})

	// Authentication Routes
	authRoutes := app.Group("api/auth")
	{
		authRoutes.Post("/register", controllers.Register)
		authRoutes.Post("/login", controllers.Login)
		authRoutes.Post("/logout", controllers.Logout)
	}

	// User Routes
	userRoutes := app.Group("api")
	{
		userRoutes.Get("/user", controllers.GetUser)
	}
}
