package tournament_management

type Vulnerable int

const (
	None Vulnerable = iota
	NS
	EW
	Both
)

func (v Vulnerable) String() string {
	switch v {
	case None:
		return "None"
	case NS:
		return "NS"
	case EW:
		return "EW"
	case Both:
		return "Both"
	default:
		return "Unknown"
	}
}

type BoardProtocol struct {
	BoardNo    int
	Vulnerable Vulnerable
	SetId      *SetId
}

func CreateBoardProtocol(tournamentId TournamentId, boardNo int, vulnerable Vulnerable) *BoardProtocol {
	if vulnerable != None && vulnerable != NS && vulnerable != EW && vulnerable != Both {
		panic("invalid vulnerable")
	}

	return &BoardProtocol{
		BoardNo:    boardNo,
		Vulnerable: vulnerable,
	}
}

type BoardProtocolDto struct {
	BoardNo    int
	Vulnerable string
}

type BoardProtocolReadRepository interface {
	FindAll(tournamentId *string) ([]BoardProtocolDto, error)
}
