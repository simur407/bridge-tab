package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "bridge-tab/internal/user/domain"

	_ "github.com/lib/pq"
)

type PostgresUserRepository struct {
	Ctx context.Context
	Tx  *sql.Tx
}

var ErrUserNotFound = errors.New("user not found")

func (r *PostgresUserRepository) Load(id *domain.UserId) (*domain.User, error) {
	var user domain.User

	row := r.Tx.QueryRowContext(r.Ctx, "SELECT id, name FROM \"user\".users WHERE id = $1", id)
	err := row.Scan(&user.State.Id, &user.State.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", ErrUserNotFound, id)
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) Save(user *domain.User) error {
	for _, event := range user.GetEvents() {
		switch event := event.(type) {
		case domain.UserRegistered:
			return r.UserRegistered(event)
		default:
			return errors.New("unknown event")
		}
	}
	return nil
}

func (r *PostgresUserRepository) UserRegistered(event domain.UserRegistered) error {
	_, err := r.Tx.ExecContext(r.Ctx, "INSERT INTO \"user\".users (id, name) VALUES ($1, $2)", event.Id, event.Name)

	return err
}

func (r *PostgresUserRepository) GetById(id string) (*domain.UserDto, error) {
	user := new(domain.UserDto)

	row := r.Tx.QueryRowContext(r.Ctx, "SELECT id, name FROM \"user\".users WHERE id = $1", id)
	err := row.Scan(&user.Id, &user.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepository) FindAll() ([]domain.UserDto, error) {
	rows, err := r.Tx.QueryContext(r.Ctx, "SELECT id, name FROM \"user\".users")
	if err != nil {
		return nil, err
	}

	var users []domain.UserDto
	for rows.Next() {
		var user domain.UserDto
		err := rows.Scan(&user.Id, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
