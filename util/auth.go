package util

import (
	"os"
	"time"

	"github.com/SOMONSOUM/go-fiber/config"
	"github.com/SOMONSOUM/go-fiber/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var db = config.SetupDatabaseConnection()
var jwtKey = []byte(os.Getenv("SECRET_KEY"))

func GenerateTokens(uuid string) (string, string) {
	claim, accessToken := GenerateAccessClaims(uuid)
	refreshToken := GenerateRefreshClaims(claim)
	return accessToken, refreshToken
}

func GenerateAccessClaims(uuid string) (*models.Claims, string) {
	t := time.Now()
	claim := &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    uuid,
			ExpiresAt: t.Add(15 * time.Minute).Unix(),
			Subject:   "access_token",
			IssuedAt:  t.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		panic(err)
	}

	return claim, tokenString
}

func GenerateRefreshClaims(cl *models.Claims) string {
	result := db.Where(&models.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer: cl.Issuer,
		},
	}).Find(&models.Claims{})

	// checking the number of refresh tokens stored.
	// If the number is higher than 3, remove all the refresh tokens and leave only new one.
	if result.RowsAffected > 3 {
		db.Where(&models.Claims{
			StandardClaims: jwt.StandardClaims{Issuer: cl.Issuer},
		}).Delete(&models.Claims{})
	}

	t := time.Now()
	refreshClaim := &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    cl.Issuer,
			ExpiresAt: cl.ExpiresAt,
			Subject:   "refresh_token",
			IssuedAt:  t.Unix(),
		},
	}

	// create a claim on database
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)

	if err != nil {
		panic(err)
	}

	return refreshTokenString
}

func SecureAuth() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		accessToken := c.Cookies("accessToken")
		claims := new(models.Claims)

		token, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if token.Valid {
			if claims.ExpiresAt < time.Now().Unix() {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   err,
					"message": "Token expired",
				})
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// this is not even a token, we should delete the cookies here
				c.ClearCookie("access_token", "refresh_token")
				return c.SendStatus(fiber.StatusForbidden)
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return c.SendStatus(fiber.StatusUnauthorized)
			} else {
				// cannot handle this token
				c.ClearCookie("access_token", "refresh_token")
				return c.SendStatus(fiber.StatusForbidden)
			}
		}

		c.Locals("id", claims.Issuer)
		return c.Next()

	}
}

func GetAuthCookies(accessToken, refreshToken string) (*fiber.Cookie, *fiber.Cookie) {
	accessCookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	}

	refreshCookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(10 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	}
	return accessCookie, refreshCookie
}
