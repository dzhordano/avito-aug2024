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

type UsersRepo struct {
	db *pgxpool.Pool
}

func NewUsersRepo(db *pgxpool.Pool) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) Create(ctx context.Context, user domain.User) error {
	const op = "repository.UsersRepo.Create"

	query, args, err := squirrel.
		Insert(usersTable).
		Columns("user_id", "email", "password_hash", "user_type").
		Values(user.ID, user.Email, user.Password, user.UserType).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {

			return fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UsersRepo) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	const op = "repository.UsersRepo.GetByCredentials"

	query, args, err := squirrel.
		Select("user_id", "email", "password_hash", "user_type").
		From(usersTable).
		Where(squirrel.And{
			squirrel.Eq{"email": email},
			squirrel.Eq{"password_hash": password},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	var user domain.User
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Email, &user.Password, &user.UserType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			
			return domain.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
