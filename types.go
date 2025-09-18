package main

import "time"

// Entity is a GTFS-RT entity with metadata and TTL semantics.
type Entity struct {
	ID         string
	Datasource string
	Kind       string // "trip_update" | "vehicle_position" | "alert"
	Message    *FeedEntity
	TTL        time.Duration
}

// Options controls conversion behavior and filters.
type Options struct {
	ETWhitelist []string
	VMWhitelist []string
	SXWhitelist []string

	CloseToNextStopPercentage int
	CloseToNextStopDistance   int

	VMGracePeriod time.Duration
}

func DefaultOptions() Options {
	return Options{
		CloseToNextStopPercentage: 95,
		CloseToNextStopDistance:   500,
		VMGracePeriod:             5 * time.Minute,
	}
}
