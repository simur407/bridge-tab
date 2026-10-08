package tournament_management

type SetId string

type Set struct {
	Id        SetId
	Label     string
	BoardNos  []int
	TeamPairs []TeamPairs
}

type SetDto struct {
	Id       string
	Label    string
	BoardNos []int
}

type SetReadRepository interface {
	FindAll(tournamentId *string) ([]SetDto, error)
}
