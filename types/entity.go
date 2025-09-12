package types

import (
	"time"

	"golang/gtfsrt"
)

// Entity is a GTFS-RT entity with metadata and TTL semantics.
type Entity struct {
	ID         string
	Datasource string
	Kind       string // "trip_update" | "vehicle_position" | "alert"
	Message    *gtfsrt.FeedEntity
	TTL        time.Duration
}
