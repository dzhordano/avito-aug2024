package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
	"github.com/dzhordano/avito-bootcamp2024/internal/dtos"
	"github.com/dzhordano/avito-bootcamp2024/internal/repository"
	"github.com/dzhordano/avito-bootcamp2024/pkg/auth"
	"github.com/google/uuid"
	"log/slog"
)

type UsersService struct {
	repo          repository.Users
	tokensManager auth.TokensManager
	log           *slog.Logger
}

func NewUsersService(repo repository.Users, tokenManager auth.TokensManager, log *slog.Logger) *UsersService {
	return &UsersService{
		repo:          repo,
		tokensManager: tokenManager,
		log:           log,
	}
}

func (s *UsersService) DummyLogin(userType string) (string, error) {
	const op = "service.Users.DummyLogin"

	log := s.log.With(
		slog.String("op", op),
		slog.String("userType", userType),
	)

	log.Info("generating auth token")

	token, err := s.tokensManager.GenerateJWT(userType)
	if err != nil {
		s.log.Error("failed to generate token: " + err.Error())

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (s *UsersService) Register(ctx context.Context, user dtos.UserRegisterInput) (string, error) {
	const op = "service.Users.Register"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", user.Email),
	)

	userId, err := uuid.NewUUID()
	if err != nil {
		s.log.Error("failed to generate user id: " + err.Error())

		return "", fmt.Errorf("%s: %w", op, err)
	}

	passwordHash := sha1.Sum([]byte(user.Password))
	user.Password = fmt.Sprintf("%x", passwordHash)

	inpUser := domain.User{
		ID:       userId,
		Email:    user.Email,
		Password: user.Password,
		UserType: user.UserType,
	}

	log.Info("registering user")

	if err = s.repo.Create(ctx, inpUser); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			s.log.Error("user already exists" + err.Error())

			return "", fmt.Errorf("%s: %w", op, domain.ErrUserAlreadyExists)
		}

		s.log.Error("failed to create user: " + err.Error())

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userId.String(), nil
}

func (s *UsersService) Login(ctx context.Context, user dtos.UserLoginInput) (string, error) {
	const op = "service.Users.Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", user.Email),
	)

	log.Info("logging in user")

	passwordHash := sha1.Sum([]byte(user.Password))
	user.Password = fmt.Sprintf("%x", passwordHash)

	respUser, err := s.repo.GetByCredentials(ctx, user.Email, user.Password)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {

			return "", fmt.Errorf("%s: %w", op, domain.ErrUserNotFound)
		}

		s.log.Error("failed to get user: " + err.Error())

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("generating auth token")

	token, err := s.tokensManager.GenerateJWT(string(respUser.UserType))
	if err != nil {
		s.log.Error("failed to generate token: " + err.Error())

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
