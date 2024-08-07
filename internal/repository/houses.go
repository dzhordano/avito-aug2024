package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HousesRepo struct {
	db *pgxpool.Pool
}

func NewHousesRepo(db *pgxpool.Pool) *HousesRepo {
	return &HousesRepo{
		db: db,
	}
}

func (r *HousesRepo) GetById(ctx context.Context, id int) ([]domain.Flat, error) {
	const op = "repository.HousesRepo.GetById"

	// Get flats to query
	query, args, err := squirrel.
		Select("flat_id").
		From(houseFlatsTable).
		Where(squirrel.Eq{"house_id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var flatsIds []int

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrHouseNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var flatId int
		if err = rows.Scan(&flatId); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		flatsIds = append(flatsIds, flatId)
	}

	// Get flats by ids
	query, args, err = squirrel.
		Select("id", "flat_number", "price", "rooms", "status").
		From(flatsTable).
		Where(squirrel.Eq{"id": flatsIds, "status": r.statusesFromUserType(ctx)}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var flats []domain.Flat
	rows, err = r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		flat := domain.Flat{}
		if err = rows.Scan(&flat.ID, &flat.FlatNumber, &flat.Price, &flat.Rooms, &flat.Status); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		flats = append(flats, flat)
	}

	return flats, nil
}

func (r *HousesRepo) statusesFromUserType(ctx context.Context) []domain.Status {
	userType := ctx.Value("user-type").(string)

	if userType == string(domain.UserTypeModerator) {
		return []domain.Status{
			domain.StatusCreated,
			domain.StatusApproved,
			domain.StatusDeclined,
			domain.StatusOnModeration,
		}
	}

	return []domain.Status{domain.StatusApproved}
}

func (r *HousesRepo) Create(ctx context.Context, house domain.House) (domain.House, error) {
	const op = "repository.HousesRepo.Create"

	query, args, err := squirrel.
		Insert(housesTable).
		Columns("address", "year", "developer", "created_at", "updated_at").
		Values(house.Address, house.Year, house.Developer, house.CreatedAt, house.UpdatedAt).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return domain.House{}, fmt.Errorf("%s: %w", op, err)
	}

	var houseId int
	err = r.db.QueryRow(ctx, query, args...).Scan(&houseId)
	if err != nil {
		return domain.House{}, fmt.Errorf("%s: %w", op, err)
	}

	house.ID = houseId

	return house, nil
}

func (r *HousesRepo) SubscribeUser(ctx context.Context, houseId int, email string) error {
	const op = "repository.HousesRepo.SubscribeUser"

	query, args, err := squirrel.
		Insert(houseSubsTable).
		Columns("house_id", "user_email").
		Values(houseId, email).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {

			return fmt.Errorf("%s: %w", op, ErrUserAlreadySubscribed)
		}

		if pgErr.Code == pgerrcode.ForeignKeyViolation {

			return fmt.Errorf("%s: %w", op, ErrUserOrHouseNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *HousesRepo) GetHouseSubscribers(ctx context.Context, houseId int) ([]string, error) {
	const op = "repository.HousesRepo.GetHouseSubscribers"

	query, args, err := squirrel.
		Select("user_email").
		From(houseSubsTable).
		Where(squirrel.Eq{"house_id": houseId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var emails []string
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		if err = rows.Scan(&email); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		emails = append(emails, email)
	}

	return emails, nil
}
