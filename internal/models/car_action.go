package models

type CarAction string

const (
	CarActionForward  CarAction = "forward"
	CarActionBackward CarAction = "backward"
	CarActionLeft     CarAction = "left"
	CarActionRight    CarAction = "right"
	CarActionStop     CarAction = "stop"
)

type CarActionRequest struct {
	Action CarAction `json:"action"`
}

func (a CarAction) IsValid() bool {
	switch a {
	case
		CarActionForward,
		CarActionBackward,
		CarActionLeft,
		CarActionRight,
		CarActionStop:
		return true
	}

	return false
}
