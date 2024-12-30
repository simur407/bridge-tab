package rounds_registration

type Round struct {
	DealNo int
	NsTeam *TeamId
	EwTeam *TeamId
	// set when round is done
	Contract    string // [1-7][CDHSN]x{0,2}|Pass
	Tricks      int    // 0-13
	Declarer    string // "N" "S" "E" "W"
	OpeningLead string // [2-9]|10|[AKQJ][CDHSN]
}

func (r *Round) IsPlayed() bool {
	return r.Contract != ""
}
