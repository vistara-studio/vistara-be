package user

import (
	"errors"

	"github.com/vistara-studio/vistara-be/pkg/cerr"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrEmailAlreadyExists = cerr.New(fiber.ErrConflict.Code, "email has been taken", errors.New("unique email constraint violation"))
	ErrUserNotFound       = cerr.New(fiber.ErrNotFound.Code, "account not found", errors.New("account not found"))
)
