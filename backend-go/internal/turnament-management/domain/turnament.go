package turnament_management

type TurnamentState struct {
	Id   string
	Name string
}

type Turnament struct {
	State   TurnamentState
	events  []any
}

// events
type TurnamentCreated struct {
	Id   string
	Name string
}

func CreateTurnament(Id string, Name string) *Turnament {
	return &Turnament{
		State:  TurnamentState{Id: Id, Name: Name},
		events: []any{TurnamentCreated{Id: Id, Name: Name}},
	}
}

// GetEvents returns the events of a turnament
func (t *Turnament) GetEvents() []any {
	return t.events
}

type TurnamentRepository interface {
	Load(Id string) (*Turnament, error)
	Save(t *Turnament) error
}

type InMemoryTurnamentRepository struct {
	turnaments []*Turnament
}

func (r *InMemoryTurnamentRepository) Load(Id string) (*Turnament, error) {
	for _, t := range r.turnaments {
		if t.State.Id == Id {
			return t, nil
		}
	}
	return nil, nil
}

func (r *InMemoryTurnamentRepository) Save(t *Turnament) error {
	for i, turnament := range r.turnaments {
		if turnament.State.Id == t.State.Id {
			r.turnaments[i] = t
			return nil
		}
	}
	r.turnaments = append(r.turnaments, t)
	return nil
}
