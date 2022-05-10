package controllers

import (
	"os"

	"github.com/SOMONSOUM/go-fiber/config"
	"github.com/SOMONSOUM/go-fiber/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

var db *gorm.DB = config.SetupDatabaseConnection()

func GetUser(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	db.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}
