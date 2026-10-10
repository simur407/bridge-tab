package main

import "testing"

func TestParsePairLines(t *testing.T) {
	pairs, err := parsePairLines("1;2\n3;4;5\n\n2:6;7\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(pairs) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(pairs))
	}
	if pairs[0].Table != nil || pairs[0].NS != 1 || pairs[0].EW != 2 {
		t.Fatalf("unexpected first pair: %+v", pairs[0])
	}
	if pairs[1].Table == nil || *pairs[1].Table != 3 || pairs[1].NS != 4 || pairs[1].EW != 5 {
		t.Fatalf("unexpected second pair: %+v", pairs[1])
	}
	if pairs[2].Table == nil || *pairs[2].Table != 2 || pairs[2].NS != 6 || pairs[2].EW != 7 {
		t.Fatalf("unexpected third pair: %+v", pairs[2])
	}
}

func TestParsePairLinesRejectsEmpty(t *testing.T) {
	if _, err := parsePairLines("  \n"); err == nil {
		t.Fatal("expected an error")
	}
}

func TestParseBoardNos(t *testing.T) {
	nums, err := parseBoardNos("1, 2,3")
	if err != nil {
		t.Fatal(err)
	}
	if len(nums) != 3 || nums[0] != 1 || nums[2] != 3 {
		t.Fatalf("unexpected boards: %v", nums)
	}
}

func TestNormalizePostedContract(t *testing.T) {
	if got := normalizePostedContract(" pass "); got != "Pass" {
		t.Fatalf("got %q", got)
	}
	if got := normalizePostedContract("3nt"); got != "3NT" {
		t.Fatalf("got %q", got)
	}
}
