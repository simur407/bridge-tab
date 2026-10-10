package score_registration_test

import (
	"testing"

	. "bridge-tab/internal/score-registration/domain"
)

func TestBridgeScorePass(t *testing.T) {
	score, err := BridgeScore("Pass", 0, "", false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score != 0 {
		t.Fatalf("expected 0, got %d", score)
	}
}

func TestBridgeScorePartscore(t *testing.T) {
	// 2NT by N, 8 tricks, NV → 70 + 50 = 120 (sheet Komplet A board 1 top)
	score, err := BridgeScore("2N", 8, "N", false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score != 120 {
		t.Fatalf("expected 120, got %d", score)
	}
}

func TestBridgeScoreMajorPartscore(t *testing.T) {
	score, err := BridgeScore("2H", 8, "N", false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score != 110 {
		t.Fatalf("expected 110, got %d", score)
	}
}

func TestBridgeScoreEWDeclarerFlipsSign(t *testing.T) {
	// 1NT by E down 1 NV → EW -50 → NS +50
	score, err := BridgeScore("1N", 6, "E", false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score != 50 {
		t.Fatalf("expected 50, got %d", score)
	}
}

func TestBridgeScoreDoubledDown(t *testing.T) {
	// 1NTx by S, 5 tricks → down 2 doubled NV = -300
	score, err := BridgeScore("1Nx", 5, "S", false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score != -300 {
		t.Fatalf("expected -300, got %d", score)
	}
}

func TestBridgeScoreGameVulnerable(t *testing.T) {
	score, err := BridgeScore("3N", 9, "N", true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 100 + 500 = 600
	if score != 600 {
		t.Fatalf("expected 600, got %d", score)
	}
}

func TestBridgeScoreSmallSlam(t *testing.T) {
	score, err := BridgeScore("6S", 12, "N", false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 180 + 300 + 500 = 980
	if score != 980 {
		t.Fatalf("expected 980, got %d", score)
	}
}

func TestBridgeScoreNTAlias(t *testing.T) {
	score, err := BridgeScore("3NT", 9, "N", false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score != 400 {
		t.Fatalf("expected 400, got %d", score)
	}
}
