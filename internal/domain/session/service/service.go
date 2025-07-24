package service

import (
	"context"

	"github.com/vistara-studio/vistara-be/internal/domain/session"
	sessionRepository "github.com/vistara-studio/vistara-be/internal/domain/session/repository"
	userRepository "github.com/vistara-studio/vistara-be/internal/domain/user/repository"
	"github.com/vistara-studio/vistara-be/pkg/jwt"
)

type authService struct {
	repository        userRepository.RepositoryItf
	sessionRepository sessionRepository.RepositoryItf
	jwt               *jwt.JWTStruct
}

type AuthServiceItf interface {
	Register(ctx context.Context, request session.RegisterRequest) (session.LoginResponse, error)
	Login(ctx context.Context, request session.LoginRequest) (session.LoginResponse, error)
}

func New(repository userRepository.RepositoryItf, sessionRepository sessionRepository.RepositoryItf, jwt *jwt.JWTStruct) AuthServiceItf {

	return &authService{
		repository:        repository,
		sessionRepository: sessionRepository,
		jwt:               jwt,
	}
}
