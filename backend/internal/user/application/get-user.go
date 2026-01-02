package user

import (
	domain "bridge-tab/internal/user/domain"
)

type GetUserCommand struct {
	Id string
}

func (cmd *GetUserCommand) Execute(repository domain.UserReadRepository) (*domain.UserDto, error) {
	u, err := repository.GetById(cmd.Id)
	if err != nil {
		return nil, err
	}

	return u, nil
}
