package umpirecall

import (
	"testing"
	"time"
)

func TestButtonFrom(t *testing.T) {
	now := time.Date(2026, 10, 20, 18, 0, 0, 0, time.UTC)
	answered := now.Add(-4 * time.Second)
	old := now.Add(-(ComingWindow + time.Second))

	state, left := ButtonFrom(nil, now)
	if state != ButtonIdle || left != 0 {
		t.Fatalf("no call: got %s %s", state, left)
	}

	state, left = ButtonFrom(&Call{}, now)
	if state != ButtonWaiting || left != 0 {
		t.Fatalf("open call: got %s %s", state, left)
	}

	state, left = ButtonFrom(&Call{AnsweredAt: &answered}, now)
	if state != ButtonComing || left != ComingWindow-4*time.Second {
		t.Fatalf("just answered: got %s %s", state, left)
	}

	state, left = ButtonFrom(&Call{AnsweredAt: &old}, now)
	if state != ButtonIdle || left != 0 {
		t.Fatalf("answered after window: got %s %s", state, left)
	}
}
