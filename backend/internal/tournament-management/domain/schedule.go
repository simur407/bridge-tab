package tournament_management

type ScheduleEntry struct {
	RoundNo  int
	TableId  TableId
	NsTeamId TeamId
	EwTeamId TeamId
	SetId    SetId
}

type ScheduleDto struct {
	RoundNo     int
	TableId     string
	TableNumber int
	NsTeamId    string
	EwTeamId    string
	NsNumber    int
	EwNumber    int
	SetId       string
	SetLabel    string
}

type ScheduleReadRepository interface {
	FindAll(tournamentId *string) ([]ScheduleDto, error)
}
