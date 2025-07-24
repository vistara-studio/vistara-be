package middleware

import (
	"errors"
	"strings"

	"github.com/vistara-studio/vistara-be/pkg/cerr"
	"github.com/vistara-studio/vistara-be/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrTokenEmpty       = cerr.New(fiber.StatusUnauthorized, "token is empty", errors.New("token can't be empty"))
	ErrInvalidTokenType = cerr.New(fiber.StatusUnauthorized, "invalid type", errors.New("invalid token type"))
)

func Authentication(jwt *jwt.JWTStruct) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authorizationHeader := ctx.GetReqHeaders()["Authorization"]
		if len(authorizationHeader) <= 0 {
			return ErrTokenEmpty
		}

		token := strings.SplitN(authorizationHeader[0], " ", 2)
		if len(token) != 2 || token[0] != "Bearer" {
			return ErrInvalidTokenType
		}

		claims, err := jwt.Decode(token[1])
		if err != nil {
			return err
		}

		ctx.Locals("user_id", claims.UserID)
		ctx.Locals("is_premium", claims.IsPremium)
		return ctx.Next()
	}
}
