package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) CreateUser(ctx context.Context, v User) (User, error) {
	var id int64

	if err := v.hashPassword(); err != nil {
		return v, fmt.Errorf("Error while hashing password: %w", err)
	}

	err := r.Db.QueryRowContext(ctx, `
		INSERT INTO users (phone, name, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`,
		v.Phone,
		v.Name,
		v.Password).Scan(&id)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code == pgerrcode.UniqueViolation {
				return v, ErrDuplicateColumn
			}
		}
		return v, fmt.Errorf("Error insert: %w", err)
	}

	v.ID = id

	return v, nil
}

func (r *Repository) UpdateUser(ctx context.Context, v User) (User, error) {
	_, err := r.Db.ExecContext(ctx, `
		UPDATE users
		SET
			phone = $1,
			name = $2
		WHERE id = $3
	`,
		v.Phone,
		v.Name,
		v.ID,
	)

	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code {
			case pgerrcode.UniqueViolation:
				return v, ErrDuplicateColumn
			case pgerrcode.NoDataFound:
				return v, ErrNotFound
			}
		}
		return v, fmt.Errorf("Error update: %w", err)
	}

	return v, nil
}

func (r *Repository) GetUser(ctx context.Context, vid int64) (User, error) {
	user := User{}

	err := r.Db.QueryRowContext(ctx, `
		SELECT id, phone, name, password, count_login, created_at, updated_at
		FROM users
		WHERE id = $1
	`, vid).Scan(
		&user.ID,
		&user.Phone,
		&user.Name,
		&user.Password,
		&user.CountLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code == pgerrcode.NoDataFound {
				return user, ErrNotFound
			}
		}
		return user, fmt.Errorf("Error on select: %w", err)
	}

	return user, nil
}

func (r *Repository) getUserForLogin(ctx context.Context, v Login) (User, error) {
	user := User{}

	err := r.Db.QueryRowContext(ctx, `
		SELECT id, phone, name, password, count_login, created_at, updated_at
		FROM users
		WHERE phone = $1
	`, v.Phone).Scan(
		&user.ID,
		&user.Phone,
		&user.Name,
		&user.Password,
		&user.CountLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code == pgerrcode.NoDataFound {
				return user, ErrNotFound
			}
		}
		return user, fmt.Errorf("Error on select for login: %w", err)
	}

	return user, nil
}

func (r *Repository) Login(ctx context.Context, v Login) (string, error) {
	user, err := r.getUserForLogin(ctx, v)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", ErrUnauthorized
		}
		return "", fmt.Errorf("Error on login: %w", err)
	}

	if !user.passwordMatch(v.Password) {
		return "", ErrUnauthorized
	}

	token, err := user.generateToken(r.JWTSecretKey)
	if err != nil {
		return "", fmt.Errorf("Error on generating token: %w", err)
	}

	_, err = r.Db.ExecContext(ctx, `
		UPDATE users
		SET
			count_login = $1
		WHERE id = $2
	`,
		user.CountLogin+1,
		user.ID,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}
