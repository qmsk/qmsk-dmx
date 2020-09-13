package api

import (
	"time"
)

type OutputStatus struct {
	Seen    time.Time
	Address string
	Port    int

	Artnet interface{} // metadata
}

type OutputID string
type Outputs map[OutputID]Output

type Output struct {
	ID        OutputID
	Universe  Universe
	Connected *time.Time

	OutputStatus
}
