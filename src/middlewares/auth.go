package middlewares

import (
	"crypto/rsa"
	"errors"
	"log"
	"strings"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/config"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthProtected(c *fiber.Ctx) error {
	return authentication(c)
}

func authentication(c *fiber.Ctx) error {
	var (
		key         *rsa.PublicKey
		bearerToken string
		jwtToken    string
		err         error
	)

	if c.Get("Authorization") != "" {
		bearerToken = c.Get("Authorization")
	} else if c.Get("authorization") != "" {
		bearerToken = c.Get("authorization")
	} else {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Split the bearer token
	split := strings.Split(bearerToken, " ")
	if len(split) < 2 {
		return c.SendStatus(utils.StatusInvalidToken)
	}

	// Set JWT Token
	jwtToken = split[1]

	// Load the public key from environment variable
	SecretKey := "-----BEGIN CERTIFICATE-----\n" + config.OAuthConfig.PublicKey + "\n-----END CERTIFICATE-----"

	// Parse a certificate
	key, err = jwt.ParseRSAPublicKeyFromPEM([]byte(SecretKey))
	if err != nil {
		// Invalid public key
		log.Println("ParseRSAPublicKeyFromPEM Error:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Verify the signature
	token, err := jwt.Parse(jwtToken, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected method: " + jwtToken.Header["alg"].(string))
		}
		return key, nil
	})
	if err != nil {
		// 401, Unexpected method token algorithm
		return c.Status(fiber.StatusUnauthorized).JSON(err)
	}

	// TODO: using claims data for get user from database
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(errors.New("invalid token"))
	}

	// TODO: Load user from database

	return nil
}
