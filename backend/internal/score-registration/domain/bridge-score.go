package score_registration

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidContract = errors.New("invalid contract")
	ErrInvalidDeclarer = errors.New("invalid declarer")
	ErrInvalidTricks   = errors.New("invalid tricks")
)

var contractPattern = regexp.MustCompile(`^([1-7])([CDHSN])(x{0,2})$`)

// BridgeScore computes the NS-perspective duplicate score for a played result.
// Vulnerability flags describe whether NS / EW are vulnerable on the board.
func BridgeScore(contract string, tricks int, declarer string, nsVulnerable, ewVulnerable bool) (int, error) {
	contract = strings.ReplaceAll(contract, "NT", "N")
	if strings.EqualFold(contract, "Pass") || contract == "Pas" {
		return 0, nil
	}

	if tricks < 0 || tricks > 13 {
		return 0, ErrInvalidTricks
	}

	matches := contractPattern.FindStringSubmatch(contract)
	if matches == nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidContract, contract)
	}

	level, _ := strconv.Atoi(matches[1])
	strain := matches[2]
	doubles := len(matches[3])

	switch declarer {
	case "N", "S", "E", "W":
	default:
		return 0, ErrInvalidDeclarer
	}

	declarerVulnerable := nsVulnerable
	if declarer == "E" || declarer == "W" {
		declarerVulnerable = ewVulnerable
	}

	needed := 6 + level
	overUnder := tricks - needed

	var score int
	if overUnder >= 0 {
		score = madeScore(level, strain, doubles, overUnder, declarerVulnerable)
	} else {
		score = -undertrickPenalty(-overUnder, doubles, declarerVulnerable)
	}

	if declarer == "E" || declarer == "W" {
		score = -score
	}
	return score, nil
}

func madeScore(level int, strain string, doubles, overtricks int, vulnerable bool) int {
	contractTricks := contractTrickScore(level, strain) * (1 << doubles)

	var overtrickPts int
	if doubles == 0 {
		per := 30
		if strain == "C" || strain == "D" {
			per = 20
		}
		overtrickPts = overtricks * per
	} else {
		per := 100 * doubles
		if vulnerable {
			per *= 2
		}
		overtrickPts = overtricks * per
	}

	var gamePart int
	if contractTricks >= 100 {
		if vulnerable {
			gamePart = 500
		} else {
			gamePart = 300
		}
	} else {
		gamePart = 50
	}

	var slam int
	if level >= 6 {
		per := 500
		if vulnerable {
			per = 750
		}
		slam = (level - 5) * per
	}

	insult := doubles * 50
	return contractTricks + overtrickPts + gamePart + slam + insult
}

func contractTrickScore(level int, strain string) int {
	switch strain {
	case "C", "D":
		return level * 20
	case "H", "S":
		return level * 30
	case "N":
		return 40 + (level-1)*30
	default:
		return 0
	}
}

func undertrickPenalty(undertricks, doubles int, vulnerable bool) int {
	if doubles == 0 {
		if vulnerable {
			return undertricks * 100
		}
		return undertricks * 50
	}

	total := 0
	for i := 1; i <= undertricks; i++ {
		var p int
		if !vulnerable {
			switch {
			case i == 1:
				p = 100
			case i <= 3:
				p = 200
			default:
				p = 300
			}
		} else {
			if i == 1 {
				p = 200
			} else {
				p = 300
			}
		}
		total += p
	}
	if doubles >= 2 {
		total *= 2
	}
	return total
}
