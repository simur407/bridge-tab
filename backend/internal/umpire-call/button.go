package umpirecall

import "time"

type ButtonState string

const (
	ButtonIdle    ButtonState = "idle"
	ButtonWaiting ButtonState = "waiting"
	ButtonComing  ButtonState = "coming"
)

// ButtonFrom reports how the calling screen should render the latest call.
// The second result is the time left until the button may be used again.
func ButtonFrom(call *Call, now time.Time) (ButtonState, time.Duration) {
	if call == nil {
		return ButtonIdle, 0
	}
	if call.AnsweredAt == nil {
		return ButtonWaiting, 0
	}
	remaining := ComingWindow - now.Sub(*call.AnsweredAt)
	if remaining > 0 {
		return ButtonComing, remaining
	}
	return ButtonIdle, 0
}
