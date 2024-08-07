package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
	"github.com/dzhordano/avito-bootcamp2024/internal/dtos"
	"github.com/dzhordano/avito-bootcamp2024/internal/repository"
	"github.com/dzhordano/avito-bootcamp2024/pkg/notifications/sender"
	"log/slog"
	"sync"
)

var ()

type FlatsService struct {
	repo          repository.Flats
	houseRepo     repository.Houses
	notifications sender.Sender
	wg            *sync.WaitGroup
	log           *slog.Logger
}

func NewFlatsService(repo repository.Flats, houseRepo repository.Houses, notifications sender.Sender, wg *sync.WaitGroup, log *slog.Logger) *FlatsService {
	return &FlatsService{
		repo:          repo,
		houseRepo:     houseRepo,
		notifications: notifications,
		wg:            wg,
		log:           log,
	}
}

func (s *FlatsService) Create(ctx context.Context, flatInp dtos.FlatCreateInput) (domain.Flat, error) {
	const op = "service.Flats.Create"

	log := s.log.With(
		slog.String("op", op),
		slog.String("flat_number", fmt.Sprint(flatInp.FlatNumber)),
		slog.Int("house_id", flatInp.HouseId),
	)

	repoFlat := domain.Flat{
		FlatNumber: flatInp.FlatNumber,
		Price:      flatInp.Price,
		Rooms:      flatInp.Rooms,
		Status:     domain.StatusCreated,
	}

	log.Info("creating flatInp")

	resp, err := s.repo.Create(ctx, flatInp.HouseId, repoFlat)
	if err != nil {
		if errors.Is(err, repository.ErrFlatAlreadyExists) {
			s.log.Error("flat already exists: " + err.Error())

			return domain.Flat{}, fmt.Errorf("%s: %w", op, ErrFlatAlreadyExists)
		}

		if errors.Is(err, repository.ErrHouseNotFound) {
			s.log.Error("house not found: " + err.Error())

			return domain.Flat{}, fmt.Errorf("%s: %w", op, ErrHouseNotFound)
		}

		s.log.Error("failed to create flatInp: " + err.Error())

		return domain.Flat{}, fmt.Errorf("%s: %w", op, err)
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		ctx := context.Background()

		var recipients []string
		subscribers, err := s.houseRepo.GetHouseSubscribers(ctx, flatInp.HouseId)
		if err != nil {
			s.log.Error("failed to get subscribers: " + err.Error())

			return
		} else {
			recipients = append(recipients, subscribers...)
		}

		log.Info("sending email")
		for _, r := range recipients {
			if err := s.notifications.SendEmail(ctx, r, fmt.Sprintf("message for house %d subscriber!", flatInp.HouseId)); err != nil {
				s.log.Error("failed to send email: " + err.Error())
			}
		}
	}()

	return resp, nil
}

func (s *FlatsService) Update(ctx context.Context, flatId int, status domain.Status) (domain.Flat, error) {
	const op = "service.Flats.Update"
	log := s.log.With(
		slog.String("op", op),
		slog.Int("flatId", flatId),
		slog.String("status", status.String()),
	)

	log.Info("switching flat status")
	// Check if flat is currently on moderation or start it.
	// Switch Moderation (naming) is due to a necessity to switch it back in case of an error.
	isBlocked, err := s.repo.SwitchModeration(ctx, flatId)
	if err != nil {
		s.log.Error("failed to switch flat status to 'moderation':" + err.Error())

		return domain.Flat{}, fmt.Errorf("%s: %w", op, err)
	}

	if isBlocked {
		s.log.Error("failed to switch flat status to 'moderation': " + ErrAlreadyModerating.Error())

		return domain.Flat{}, fmt.Errorf("%s: %w", op, ErrAlreadyModerating)
	}

	log.Info("updating flat status")

	resp, err := s.repo.Update(ctx, flatId, status.String())
	if err != nil {
		err = s.repo.SwitchModerationBackTo(context.Background(), flatId, status.String())
		if err != nil {
			s.log.Error("failed to update and switch flat status to back: " + err.Error())

			return domain.Flat{}, fmt.Errorf("%s: %w", op, err)
		}

		if ctx.Err() != nil {
			s.log.Error("failed to update flat status due to context error: " + ctx.Err().Error())

			return domain.Flat{}, fmt.Errorf("%s: %w", op, ctx.Err())
		}

		if errors.Is(err, repository.ErrFlatNotFound) {
			s.log.Error("flat not found: " + err.Error())

			return domain.Flat{}, fmt.Errorf("%s: %w", op, ErrFlatNotFound)
		}

		s.log.Error("failed to update flat status: " + err.Error())

		return domain.Flat{}, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}
