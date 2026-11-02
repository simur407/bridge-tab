package tournament_management

import (
	"errors"
	"strings"
)

var (
	ErrInvalidHandFormat    = errors.New("invalid hand format")
	ErrInvalidHandRank      = errors.New("invalid hand rank")
	ErrDuplicateCardInHand  = errors.New("duplicate card in hand")
	ErrHandWrongLength      = errors.New("hand must contain 13 cards")
	ErrDealDuplicateCard    = errors.New("card appears in more than one hand")
	ErrDealIncomplete       = errors.New("deal must use all 52 cards")
	ErrBoardHandsAlreadySet = errors.New("board hands already recorded")
)

const HAND_RANKS = "AKQJT98765432"

type BoardHands struct {
	North string
	East  string
	South string
	West  string
}

func (h BoardHands) Complete() bool {
	return h.North != "" && h.East != "" && h.South != "" && h.West != ""
}

func NormalizeAndValidateDeal(north, east, south, west string) (BoardHands, error) {
	n, nCards, err := normalizeHand(north)
	if err != nil {
		return BoardHands{}, err
	}
	e, eCards, err := normalizeHand(east)
	if err != nil {
		return BoardHands{}, err
	}
	s, sCards, err := normalizeHand(south)
	if err != nil {
		return BoardHands{}, err
	}
	w, wCards, err := normalizeHand(west)
	if err != nil {
		return BoardHands{}, err
	}

	// we want a set of all cards, struct{} is for memory optimization
	seen := make(map[int]struct{}, 52)

	for _, card := range append(append(append(nCards, eCards...), sCards...), wCards...) {
		// check if already seen
		_, ok := seen[card]
		if ok {
			return BoardHands{}, ErrDealDuplicateCard
		}
		seen[card] = struct{}{}
	}
	if len(seen) != 52 {
		return BoardHands{}, ErrDealIncomplete
	}

	return BoardHands{North: n, East: e, South: s, West: w}, nil
}

// transforming raw hand to normalized hand.
// It sorts ranks, removes duplicates and converts 10 to T.
// Then it changes cards to numbers, so we can check for duplicates.
func normalizeHand(raw string) (string, []int, error) {
	parts := strings.Split(raw, ".")
	if len(parts) != 4 {
		return "", nil, ErrInvalidHandFormat
	}

	cards := make([]int, 0, 13)
	normalized := make([]string, 4)
	for suit, part := range parts {
		ranks := make([]byte, 0, len(part))
		usedInSuit := make(map[int]struct{}, len(part))
		for i := 0; i < len(part); i++ {
			rank := part[i]
			if rank == '1' && i+1 < len(part) && part[i+1] == '0' {
				rank = 'T'
				i++
			}
			rankIndex := strings.IndexByte(HAND_RANKS, rank)
			if rankIndex < 0 {
				return "", nil, ErrInvalidHandRank
			}
			if _, ok := usedInSuit[rankIndex]; ok {
				return "", nil, ErrDuplicateCardInHand
			}
			usedInSuit[rankIndex] = struct{}{}
			ranks = append(ranks, HAND_RANKS[rankIndex])
			cards = append(cards, suit*13+rankIndex)
		}
		sortRanksHighToLow(ranks)
		normalized[suit] = string(ranks)
	}
	if len(cards) != 13 {
		return "", nil, ErrHandWrongLength
	}
	return strings.Join(normalized, "."), cards, nil
}

func sortRanksHighToLow(ranks []byte) {
	for i := 1; i < len(ranks); i++ {
		j := i
		for j > 0 && rankOrder(ranks[j-1]) > rankOrder(ranks[j]) {
			ranks[j-1], ranks[j] = ranks[j], ranks[j-1]
			j--
		}
	}
}

func rankOrder(rank byte) int {
	return strings.IndexByte(HAND_RANKS, rank)
}
