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
	North      string
	East       string
	South      string
	West       string
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

func (bp *BoardProtocol) Hands() BoardHands {
	return BoardHands{
		North: bp.North,
		East:  bp.East,
		South: bp.South,
		West:  bp.West,
	}
}

func (bp *BoardProtocol) HandsComplete() bool {
	return bp.Hands().Complete()
}

type BoardProtocolDto struct {
	BoardNo    int
	Vulnerable string
	North      string
	East       string
	South      string
	West       string
}

func (d BoardProtocolDto) HandsComplete() bool {
	return d.North != "" && d.East != "" && d.South != "" && d.West != ""
}

type BoardProtocolReadRepository interface {
	FindAll(tournamentId *string) ([]BoardProtocolDto, error)
	FindByBoardNo(tournamentId *string, boardNo int) (*BoardProtocolDto, error)
}
