package user

import (
	"errors"

	domain "bridge-tab/internal/user/domain"
	infra "bridge-tab/internal/user/infrastructure"
)

type RegisterUserCommand struct {
	Id   string
	Name string
}

func (cmd *RegisterUserCommand) Execute(repository domain.UserRepository) error {
	userId := domain.UserId(cmd.Id)

	u, err := repository.Load(&userId)
	if err != nil {
		if !errors.Is(err, infra.ErrUserNotFound) {
			return err
		}
	}

	if u != nil {
		return errors.New("user already exists")
	}

	user := domain.RegisterUser(&userId, cmd.Name)
	return repository.Save(user)
}
