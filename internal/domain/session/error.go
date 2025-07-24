package session

import (
	"errors"

	"github.com/vistara-studio/vistara-be/pkg/cerr"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrSessionNotFound = cerr.New(fiber.ErrNotFound.Code, "session not found", errors.New("session not found"))
)
