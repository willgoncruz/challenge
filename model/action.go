package model

type ActionType string

// Action names
const (
	Place   ActionType = "place"
	Move    ActionType = "move"
	Pickup  ActionType = "pickup"
	Discard ActionType = "discard"
)

// Action is a json-friendly representation of an action.
type Action struct {
	Timestamp int64      `json:"timestamp"` // unix timestamp in microseconds
	ID        string     `json:"id"`        // order id
	Action    ActionType `json:"action"`    // place, move, pickup or discard (ActionType)
}
