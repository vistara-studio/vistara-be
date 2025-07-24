package local

import (
	"errors"

	"github.com/vistara-studio/vistara-be/pkg/cerr"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrLBNotFound = cerr.New(fiber.ErrNotFound.Code, "local business not found", errors.New("account not found"))
)
