package play

import "testing"

func TestValidateScoreAllowsZeroTricks(t *testing.T) {
	if err := validateScore("1C", 0, "N", "AS"); err != nil {
		t.Fatalf("1C taking no tricks: %v", err)
	}
}

func TestValidateScoreRejectsTricksOutsideRange(t *testing.T) {
	if err := validateScore("1C", -1, "N", "AS"); err == nil {
		t.Fatal("expected negative tricks to be rejected")
	}
	if err := validateScore("1C", 14, "N", "AS"); err == nil {
		t.Fatal("expected 14 tricks to be rejected")
	}
}
