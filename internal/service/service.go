package service

import (
	"context"
	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
	"github.com/dzhordano/avito-bootcamp2024/internal/dtos"
	"github.com/dzhordano/avito-bootcamp2024/internal/repository"
	"github.com/dzhordano/avito-bootcamp2024/pkg/auth"
	"github.com/dzhordano/avito-bootcamp2024/pkg/notifications/sender"
	"log/slog"
	"sync"
)

//go:generate mockgen -source=service.go -destination=mocks/mocks.go -package=service

type Houses interface {
	GetById(ctx context.Context, id int) ([]domain.Flat, error)
	Create(ctx context.Context, house dtos.HouseCreateInput) (domain.House, error)

	Subscribe(ctx context.Context, houseId int, email string) error
}

type Flats interface {
	Create(ctx context.Context, flat dtos.FlatCreateInput) (domain.Flat, error)

	Update(ctx context.Context, flatId int, status domain.Status) (domain.Flat, error)
}

type Users interface {
	DummyLogin(userType string) (string, error)
	Register(ctx context.Context, user dtos.UserRegisterInput) (string, error)
	Login(ctx context.Context, user dtos.UserLoginInput) (string, error)
}

type Services struct {
	Houses Houses
	Flats  Flats
	Users  Users
}

type Deps struct {
	Repos         *repository.Repository
	TokensManager auth.TokensManager
	Notifications sender.Sender
	WaitGroup     *sync.WaitGroup
	Logger        *slog.Logger
}

func New(deps Deps) *Services {
	users := NewUsersService(deps.Repos.Users, deps.TokensManager, deps.Logger)
	flats := NewFlatsService(deps.Repos.Flats, deps.Repos.Houses, deps.Notifications, deps.WaitGroup, deps.Logger)
	houses := NewHousesService(deps.Repos.Houses, deps.Logger)

	return &Services{
		Users:  users,
		Flats:  flats,
		Houses: houses,
	}
}
