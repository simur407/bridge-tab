package tournament_management_test

import (
	"testing"

	. "bridge-tab/internal/tournament-management/domain"
)

const (
	allSpades   = "AKQJT98765432..."
	allHearts   = ".AKQJT98765432.."
	allDiamonds = "..AKQJT98765432."
	allClubs    = "...AKQJT98765432"
)

func TestNormalizeAndValidateDeal(t *testing.T) {
	hands, err := NormalizeAndValidateDeal(allSpades, allHearts, allDiamonds, allClubs)
	assertNoError(t, err)
	if hands.North != allSpades {
		t.Fatalf("north: got %q", hands.North)
	}
	if hands.East != allHearts {
		t.Fatalf("east: got %q", hands.East)
	}
}

func TestNormalizeAndValidateDealSortsRanks(t *testing.T) {
	hands, err := NormalizeAndValidateDeal(
		"23456789TJQKA...",
		allHearts,
		allDiamonds,
		allClubs,
	)
	assertNoError(t, err)
	if hands.North != allSpades {
		t.Fatalf("expected sorted north, got %q", hands.North)
	}
}

func TestNormalizeAndValidateDealAcceptsTenAs10(t *testing.T) {
	hands, err := NormalizeAndValidateDeal(
		"AKQJ1098765432...",
		allHearts,
		allDiamonds,
		allClubs,
	)
	assertNoError(t, err)
	if hands.North != allSpades {
		t.Fatalf("expected T for ten, got %q", hands.North)
	}
}

func TestNormalizeAndValidateDealRejectsBadFormat(t *testing.T) {
	_, err := NormalizeAndValidateDeal("AKQ", allHearts, allDiamonds, allClubs)
	assertError(t, err, ErrInvalidHandFormat)
}

func TestNormalizeAndValidateDealRejectsWrongCount(t *testing.T) {
	_, err := NormalizeAndValidateDeal(
		"AKQJT9876543...",
		allHearts,
		allDiamonds,
		allClubs,
	)
	assertError(t, err, ErrHandWrongLength)
}

func TestNormalizeAndValidateDealRejectsDuplicatesAcrossHands(t *testing.T) {
	_, err := NormalizeAndValidateDeal(
		allSpades,
		"A.KQJT98765432..",
		allDiamonds,
		allClubs,
	)
	assertError(t, err, ErrDealDuplicateCard)
}

func TestNormalizeAndValidateDealRejectsBadRank(t *testing.T) {
	_, err := NormalizeAndValidateDeal(
		"AKQJT9876543X...",
		allHearts,
		allDiamonds,
		allClubs,
	)
	assertError(t, err, ErrInvalidHandRank)
}

func TestNormalizeAndValidateDealVoidSuit(t *testing.T) {
	hands, err := NormalizeAndValidateDeal(
		allSpades,
		allHearts,
		"..AKQJT98.765432",
		"..765432.AKQJT98",
	)
	assertNoError(t, err)
	if hands.South != "..AKQJT98.765432" {
		t.Fatalf("south: got %q", hands.South)
	}
	if hands.West != "..765432.AKQJT98" {
		t.Fatalf("west: got %q", hands.West)
	}
}
