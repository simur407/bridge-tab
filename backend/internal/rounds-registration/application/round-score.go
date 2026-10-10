package rounds_registration

import (
	"errors"
	"regexp"
	"strings"
)

func normalizeContract(contract string) string {
	return strings.ReplaceAll(contract, "NT", "N")
}

func validateRoundInput(gameSessionId string, playerId string, versusTeamNumber int, dealNo int, contract string, tricks int, declarer string, openingLead string) error {
	if gameSessionId == "" {
		return errors.New("game session id is empty")
	}
	if playerId == "" {
		return errors.New("player id is empty")
	}
	if versusTeamNumber < 1 {
		return errors.New("versus team number is empty")
	}
	if dealNo == 0 {
		return errors.New("deal no is empty")
	}
	return validateRoundScore(contract, tricks, declarer, openingLead)
}

func validateRoundScore(contract string, tricks int, declarer string, openingLead string) error {
	if contract == "" {
		return errors.New("contract is empty")
	}
	if contract != "Pass" && tricks == 0 {
		return errors.New("tricks is empty")
	}
	if contract != "Pass" && declarer == "" {
		return errors.New("declarer is empty")
	}
	if contract != "Pass" && openingLead == "" {
		return errors.New("opening lead is empty")
	}

	match, err := regexp.MatchString("[1-7]([CDHS]|NT?)x{0,2}|Pass", contract)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("invalid contract")
	}

	if contract != "Pass" && (tricks < 0 || tricks > 13) {
		return errors.New("invalid tricks")
	}

	if contract != "Pass" && declarer != "N" && declarer != "E" && declarer != "S" && declarer != "W" {
		return errors.New("invalid declarer")
	}

	if contract != "Pass" {
		match, err = regexp.MatchString("([2-9]|10|[AKQJ])[CDHS]", openingLead)
		if err != nil {
			return err
		}
		if !match {
			return errors.New("invalid opening lead")
		}
	}

	return nil
}
