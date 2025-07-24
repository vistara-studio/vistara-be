package service

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/vistara-studio/vistara-be/internal/domain/session"
	"github.com/vistara-studio/vistara-be/internal/domain/user"
	"github.com/vistara-studio/vistara-be/pkg/bcrypt"
	"github.com/google/uuid"
)

func (s *authService) Register(ctx context.Context, request session.RegisterRequest) (session.LoginResponse, error) {
	userRepository, err := s.repository.NewClient(false)
	if err != nil {
		return session.LoginResponse{}, err
	}

	hashedPassword, err := bcrypt.EncryptPassword(request.Password)
	if err != nil {
		return session.LoginResponse{}, err
	}

	userID, err := uuid.NewV7()
	if err != nil {
		return session.LoginResponse{}, err
	}

	user := user.Table{
		ID:           userID,
		FullName:     request.FullName,
		Email:        request.Email,
		Password:     hashedPassword,
		AuthProvider: user.AuthProviderEmail,
		PhotoUrl:     "https://htnqkjejgcovkehhtqjw.supabase.co/storage/v1/object/public/hackfest-uhuy//photo_profile.jpg",
	}

	err = userRepository.CreateUser(ctx, user)
	if err != nil {
		return session.LoginResponse{}, err
	}

	return s.Login(ctx, session.LoginRequest{
		Email:    request.Email,
		Password: request.Password,
	})
}

func (s *authService) Login(ctx context.Context, request session.LoginRequest) (session.LoginResponse, error) {
	userRepository, err := s.repository.NewClient(false)
	if err != nil {
		return session.LoginResponse{}, err
	}

	user := &user.Table{
		Email: request.Email,
	}

	if err := userRepository.GetAccountByEmail(ctx, user); err != nil {
		return session.LoginResponse{}, err
	}

	if err := bcrypt.ComparePassword(user.Password, request.Password); err != nil {
		return session.LoginResponse{}, err
	}

	token, err := s.jwt.Encode(user)
	if err != nil {
		return session.LoginResponse{}, err
	}

	sessionRepository, err := s.sessionRepository.NewClient(true)
	if err != nil {
		return session.LoginResponse{}, err
	}

	defer func() {
		if err != nil {
			errTx := sessionRepository.Rollback()
			if errTx != nil {
				err = errTx
			}
		}
	}()

	sessionID, err := uuid.NewV7()
	if err != nil {
		return session.LoginResponse{}, err
	}

	newSession := session.Table{
		ID:     sessionID,
		UserID: user.ID,
	}

	err = sessionRepository.CreateSession(ctx, newSession)
	if err != nil {
		return session.LoginResponse{}, err
	}

	// if s > 3; then delete 1
	sessions := new([]session.Table)
	err = sessionRepository.GetSessionByUserID(ctx, user, sessions)
	if err != nil {
		fmt.Println("err: " + err.Error())
		return session.LoginResponse{}, err
	}

	if len(*sessions) > 3 {
		err := sessionRepository.DeleteOldestSessionByUserID(ctx, newSession)
		if err != nil {
			return session.LoginResponse{}, err
		}
	}

	if err := sessionRepository.Commit(); err != nil {
		return session.LoginResponse{}, err
	}

	return session.LoginResponse{
		AccessToken:  token,
		RefreshToken: base64.StdEncoding.EncodeToString([]byte(newSession.ID.String())),
	}, nil
}
