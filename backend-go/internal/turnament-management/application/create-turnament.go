package turnament_management

import (
	"errors"

	domain "bridge-tab/internal/turnament-management/domain"
)

// CreateTurnamentCommand represents a command to create a turnament
type CreateTurnamentCommand struct {
	Id   string
	Name string
}

// Execute executes the command
func (c *CreateTurnamentCommand) Execute(repo domain.TurnamentRepository) error {
	// Validate input
	if err := validate(c); err != nil {
		return err
	}

	// Load the turnament from the repository
	t, err := repo.Load(c.Id)
	if err != nil {
		return err
	}

	// If the turnament already exists, return an error
	if t != nil {
		return errors.New("turnament already exists")
	}

	// Create a new turnament
	t = domain.CreateTurnament(c.Id, c.Name)

	// Save the turnament to the repository
	return repo.Save(t)
}

func validate(c *CreateTurnamentCommand) error {
	if c.Id == "" || c.Name == "" {
		return errors.New("invalid input")
	}
	return nil
}
