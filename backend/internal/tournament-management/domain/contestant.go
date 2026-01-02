package tournament_management

type ContestantId string

type Contestant struct {
	Id   ContestantId
	Team *Team
}

type ContestantDto struct {
	Id ContestantId
}
