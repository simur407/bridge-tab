package turnament_management

import (
	"database/sql"

	domain "bridge-tab/internal/turnament-management/domain"
)

type PostgresTurnamentRepository struct {
 db *sql.DB
}

func (r *PostgresTurnamentRepository) Load(Id string) (*domain.Turnament, error) {
	var turnament domain.Turnament
	row := r.db.QueryRow("SELECT id, name FROM turnament_management.turnament WHERE id = ?", Id)
	err := row.Scan(&turnament.State.Id, &turnament.State.Name);

	if err != nil {
		return nil, err
	}
	
	return &turnament, nil
}


func (r *PostgresTurnamentRepository) Save(t *domain.Turnament) error {
	for _, event := range t.GetEvents() {
		switch event.(type) {
		case domain.TurnamentCreated:

		default: 
		// todo
		}
	}
}
