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
	"time"
)

type FlatsRepo struct {
	db *pgxpool.Pool
}

func NewFlatsRepo(db *pgxpool.Pool) *FlatsRepo {
	return &FlatsRepo{
		db: db,
	}
}

func (r *FlatsRepo) Create(ctx context.Context, houseId int, flat domain.Flat) (domain.Flat, error) {
	const op = "repository.Flats.Create"

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Flat{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = fmt.Errorf("%s: %w", op, rollbackErr)
			}
		}
	}()

	// Create flat
	query, args, err := squirrel.
		Insert(flatsTable).
		Columns("flat_number", "price", "rooms", "status").
		Values(flat.FlatNumber, flat.Price, flat.Rooms, flat.Status).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return domain.Flat{}, fmt.Errorf("%s: %w", op, err)
	}

	var flatId int
	err = tx.QueryRow(ctx, query, args...).Scan(&flatId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.Flat{}, fmt.Errorf("%s: %w", op, ErrFlatAlreadyExists)
			}
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Flat{}, fmt.Errorf("%s: %w", op, ErrHouseNotFound)
		}

		return domain.Flat{}, fmt.Errorf("%s: %w", op, err)
	}

	// Add flat to house'tests flats
	query, args, err = squirrel.
		Insert(houseFlatsTable).
		Columns("flat_number", "flat_id", "house_id").
		Values(flat.FlatNumber, flatId, houseId).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return domain.Flat{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.Flat{}, fmt.Errorf("%s: %w", op, ErrFlatAlreadyExists)
			}
		}

		if pgErr.Code == pgerrcode.ForeignKeyViolation {
			return domain.Flat{}, fmt.Errorf("%s: %w", op, ErrHouseNotFound)
		}

		return domain.Flat{}, fmt.Errorf("%s: %w", op, err)
	}

	// Update house
	query, args, err = squirrel.
		Update(housesTable).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": houseId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return domain.Flat{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {

		return domain.Flat{}, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.Flat{}, fmt.Errorf("%s: %w", op, err)
	}

	// Add id to response
	flat.ID = flatId

	return flat, nil
}

func (r *FlatsRepo) Update(ctx context.Context, flatId int, status string) (domain.Flat, error) {
	const op = "repository.Flats.Update"

	// To simulate slow network
	// time.Sleep(5 * time.Second)

	query, args, err := squirrel.
		Update(flatsTable).
		Set("status", status).
		Where(squirrel.Eq{"id": flatId}).
		Suffix("RETURNING *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return domain.Flat{}, fmt.Errorf("%s: %w", op, err)
	}

	// Update flat and return it.
	var flat domain.Flat
	err = r.db.QueryRow(ctx, query, args...).Scan(&flat.ID, &flat.FlatNumber, &flat.Price, &flat.Rooms, &flat.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.NoDataFound {
				return domain.Flat{}, fmt.Errorf("%s: %w", op, ErrFlatNotFound)
			}
		}

		return domain.Flat{}, fmt.Errorf("%s: %w", op, err)
	}

	return flat, nil
}

func (r *FlatsRepo) SwitchModeration(ctx context.Context, flatId int) (bool, error) {
	const op = "repository.Flats.SwitchModeration"

	query, args, err := squirrel.
		Select("status").
		From(flatsTable).
		Where(squirrel.Eq{"id": flatId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return true, fmt.Errorf("%s: %w", op, err)
	}

	// Get flat status and check if it'tests on moderation
	var status string
	err = r.db.QueryRow(ctx, query, args...).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true, fmt.Errorf("%s: %s", op, ErrFlatNotFound)
		}

		return true, fmt.Errorf("%s: %w", op, err)
	}

	if domain.Status(status) == domain.StatusOnModeration {
		return true, fmt.Errorf("%s: %s", op, ErrFlatOnModeration)
	}

	query, args, err = squirrel.
		Update(flatsTable).
		Set("status", domain.StatusOnModeration).
		Where(squirrel.Eq{"id": flatId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return true, fmt.Errorf("%s: %w", op, err)
	}

	// Update flat status to 'moderating'.
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {

		return true, fmt.Errorf("%s: %w", op, err)
	}

	return false, nil
}

func (r *FlatsRepo) SwitchModerationBackTo(ctx context.Context, flatId int, status string) error {
	const op = "repository.Flats.SwitchModerationBackTo"

	query, args, err := squirrel.
		Update(flatsTable).
		Set("status", status).
		Where(squirrel.Eq{"id": flatId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Simply change flat status back to 'status'.
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
