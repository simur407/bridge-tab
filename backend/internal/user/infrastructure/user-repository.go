package user

import (
	"database/sql"
	"errors"
	"fmt"

	domain "bridge-tab/internal/user/domain"

	_ "github.com/lib/pq"
)

type PostgresUserRepository struct {
	Db *sql.DB
}

var ErrUserNotFound = errors.New("user not found")

func (r *PostgresUserRepository) Load(id *domain.UserId) (*domain.User, error) {
	var user domain.User

	row := r.Db.QueryRow("SELECT id, name FROM user WHERE id = $1", id)
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
	_, err := r.Db.Exec("INSERT INTO user (id, name) VALUES ($1, $2)", event.Id, event.Name)

	return err
}
