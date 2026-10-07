package tournament_management

type TableId string

type TableState struct {
	Id           TableId
	TournamentId TournamentId
	Number       int
}

type Table struct {
	State TableState
}

func CreateTable(id TableId, tournamentId TournamentId, number int) *Table {
	return &Table{
		State: TableState{Id: id, TournamentId: tournamentId, Number: number},
	}
}

type TableDto struct {
	Id           string
	TournamentId string
	Number       int
}

type TableReadRepository interface {
	FindAll(tournamentId *string) ([]TableDto, error)
	FindById(tournamentId *string, tableId *string) (*TableDto, error)
	FindByNumber(tournamentId *string, number int) (*TableDto, error)
}
