package tournament_management

type TeamPairs struct {
	NS TeamId
	EW TeamId
}

type Vulnerable int

const (
	None Vulnerable = iota
	NS
	EW
	Both
)

type BoardProtocol struct {
	BoardNo    int
	Vulnerable Vulnerable
	TeamPairs  []TeamPairs
}

func CreateBoardProtocol(tournamentId TournamentId, boardNo int, vulnerable Vulnerable, teamPairs []TeamPairs) *BoardProtocol {
	if vulnerable != None && vulnerable != NS && vulnerable != EW && vulnerable != Both {
		panic("invalid vulnerable")
	}

	return &BoardProtocol{
		BoardNo:    boardNo,
		Vulnerable: vulnerable,
		TeamPairs:  teamPairs,
	}
}

type TeamPairsDto struct {
	NS string
	EW string
}

type BoardProtocolDto struct {
	BoardNo    int
	Vulnerable string
	TeamPairs  []TeamPairsDto
}

type BoardProtocolReadRepository interface {
	FindAll(tournamentId *string) ([]BoardProtocolDto, error)
}
