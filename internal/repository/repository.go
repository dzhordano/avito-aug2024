package repository

import (
	"context"
	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	usersTable      = "users"
	housesTable     = "houses"
	flatsTable      = "flats"
	houseFlatsTable = "house_flats"
	houseSubsTable  = "house_subscriptions"
)

type Repository struct {
	Houses Houses
	Flats  Flats
	Users  Users
}

type Deps struct {
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		Houses: NewHousesRepo(db),
		Flats:  NewFlatsRepo(db),
		Users:  NewUsersRepo(db),
	}
}

type Houses interface {
	GetById(ctx context.Context, id int) ([]domain.Flat, error)
	Create(ctx context.Context, house domain.House) (domain.House, error)

	SubscribeUser(ctx context.Context, houseId int, email string) error
	GetHouseSubscribers(ctx context.Context, houseId int) ([]string, error)
}

type Flats interface {
	Create(ctx context.Context, houseId int, flat domain.Flat) (domain.Flat, error)

	Update(ctx context.Context, flatId int, status string) (domain.Flat, error)
	SwitchModeration(ctx context.Context, flatId int) (bool, error)
	SwitchModerationBackTo(ctx context.Context, flatId int, status string) error
}

type Users interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
}
