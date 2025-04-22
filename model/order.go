package model

import "time"

type Temperature string

const (
	Cold Temperature = "cold"
	Hot  Temperature = "hot"
	Room Temperature = "room"
)

// Order is a json-friendly representation of an order.
type Order struct {
	ID        string      `json:"id"`        // order id
	Name      string      `json:"name"`      // food name
	Temp      Temperature `json:"temp"`      // ideal temperature
	Freshness int         `json:"freshness"` // freshness in seconds
	TTL       int64       `json:"ttl"`       // for control of storage duration
}

func (o Order) FreshnessInSecondsByStorage(s Storage) int {
	if s.IsIdealTemp(o.Temp) {
		return o.Freshness
	}

	return o.Freshness / 2 // Cut freshness in half if not ideal storage
}

func (o Order) FillTTL(s Storage) {
	if o.TTL != 0 { // Already filled time to discard
		return
	}

	// Calculate the time to discard by freshness
	freshnessTTL := time.Duration(o.FreshnessInSecondsByStorage(s)) * time.Second
	o.TTL = time.Now().UnixMicro() + freshnessTTL.Microseconds()
}
