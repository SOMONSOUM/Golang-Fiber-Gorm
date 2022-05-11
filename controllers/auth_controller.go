package controllers

import (
	"math/rand"
	"time"

	"github.com/SOMONSOUM/go-fiber/models"
	"github.com/SOMONSOUM/go-fiber/util"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(&user); err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"input": "Please review your input",
		})
	}

	// validate if the email, username and password are in correct format
	errors := util.ValidateRegister(user)
	if errors.Err {
		return c.JSON(errors)
	}

	if count := db.Where(&models.User{Username: user.Email}).First(new(models.User)).RowsAffected; count > 0 {
		errors.Err, errors.Email = true, "Email is already registered"
	}

	if count := db.Where(&models.User{Username: user.Username}).First(new(models.User)).RowsAffected; count > 0 {
		errors.Err, errors.Email = true, "Username is already registered"
	}

	if errors.Err {
		return c.JSON(errors)
	}

	password := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(
		password, rand.Intn(bcrypt.MaxCost-bcrypt.MinCost)+bcrypt.MinCost,
	)

	if err != nil {
		panic(err)
	}

	user.Password = string(hashedPassword)
	if err := db.Create(&user).Error; err != nil {
		return c.JSON(fiber.Map{
			"error":   err,
			"message": "Something went wrong, please try again later 😕",
		})

	}

	// setting up the authorization cookies
	accessToken, refreshToken := util.GenerateTokens(user.UUID.String())
	accessCookie, refreshCookie := util.GetAuthCookies(accessToken, refreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	input := new(LoginInput)

	if err := c.BodyParser(input); err != nil {
		return c.JSON(fiber.Map{"error": true, "input": "Please review your input"})
	}

	// check if a user exists
	user := new(models.User)
	if res := db.Where(
		&models.User{Email: input.Identity}).Or(
		&models.User{Username: input.Identity},
	).First(&user); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "general": "Invalid Credentials."})
	}

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.JSON(fiber.Map{"error": true, "message": "Invalid Credentials."})
	}

	// setting up the authorization cookies
	accessToken, refreshToken := util.GenerateTokens(user.UUID.String())
	accessCookie, refreshCookie := util.GetAuthCookies(accessToken, refreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Success!",
	})
}
