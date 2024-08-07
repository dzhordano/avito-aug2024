package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
	"github.com/dzhordano/avito-bootcamp2024/internal/dtos"
	"github.com/dzhordano/avito-bootcamp2024/internal/repository"
	"log/slog"
	"time"
)

type HousesService struct {
	repo repository.Houses

	log *slog.Logger
}

func NewHousesService(repo repository.Houses, log *slog.Logger) *HousesService {
	return &HousesService{
		repo: repo,
		log:  log,
	}
}

func (s *HousesService) GetById(ctx context.Context, id int) ([]domain.Flat, error) {
	const op = "service.Houses.GetById"
	log := s.log.With(
		slog.String("op", op),
		slog.Int("house_id", id),
	)

	log.Info("collecting house flats")

	resp, err := s.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrHouseNotFound) {
			s.log.Error("house not found: " + err.Error())

			return nil, fmt.Errorf("%s: %w", op, ErrHouseNotFound)
		}

		s.log.Error("failed to get house flats")

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (s *HousesService) Create(ctx context.Context, house dtos.HouseCreateInput) (domain.House, error) {
	const op = "service.Houses.Create"

	log := s.log.With(
		slog.String("op", op),
		slog.String("address", house.Address),
	)

	houseRepo := domain.House{
		Address:   house.Address,
		Year:      house.Year,
		Developer: house.Developer,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	log.Info("creating house")

	resp, err := s.repo.Create(ctx, houseRepo)
	if err != nil {

		if errors.Is(err, repository.ErrHouseAlreadyExists) {
			s.log.Error("house already exists: " + err.Error())

			return domain.House{}, fmt.Errorf("%s: %w", op, ErrHouseAlreadyExists)
		}

		s.log.Error("failed to create house: " + err.Error())

		return domain.House{}, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (s *HousesService) Subscribe(ctx context.Context, houseId int, email string) error {
	const op = "service.Houses.Subscribe"

	log := s.log.With(
		slog.String("op", op),
		slog.Int("house_id", houseId),
		slog.String("email", email),
	)

	log.Info("subscribing user to house")

	if err := s.repo.SubscribeUser(ctx, houseId, email); err != nil {

		if errors.Is(err, repository.ErrUserAlreadySubscribed) {
			s.log.Error("user already subscribed: " + err.Error())

			return fmt.Errorf("%s: %w", op, ErrUserAlreadySubscribed)
		}

		if errors.Is(err, repository.ErrHouseNotFound) {
			s.log.Error("house not found: " + err.Error())

			return fmt.Errorf("%s: %w", op, ErrHouseNotFound)
		}

		s.log.Error("failed to subscribe user: " + err.Error())

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
