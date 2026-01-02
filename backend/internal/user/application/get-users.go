package user

import (
	domain "bridge-tab/internal/user/domain"
)

type GetUsersCommand struct {
}

func (cmd *GetUsersCommand) Execute(repository domain.UserReadRepository) ([]domain.UserDto, error) {
	u, err := repository.FindAll()
	if err != nil {
		return nil, err
	}

	return u, nil
}
