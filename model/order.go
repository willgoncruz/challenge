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
	ID        string         `json:"id"`        // order id
	Name      string         `json:"name"`      // food name
	Temp      Temperature    `json:"temp"`      // ideal temperature
	Freshness int            `json:"freshness"` // freshness in seconds
	TTL       *time.Duration `json:"ttl"`       // for control of storage duration
}
