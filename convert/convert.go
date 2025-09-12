package convert

// Package convert exposes the public API for converting SIRI (XML) to GTFS-Realtime.
// Implementation details live in implementation.go.

import (
	"golang/gtfsrt"
	"golang/siri"
	"golang/types"
)

// ConvertSIRI converts a parsed SIRI ServiceDelivery into GTFS-RT entities.
// The ServiceDelivery struct is defined in the siri package and populated via the siri/xml decoder.
func ConvertSIRI(sd *siri.ServiceDelivery, opts types.Options) ([]types.Entity, error) {
	return convertSIRI(sd, opts)
}

// BuildFeedMessage produces a single FeedMessage from entities.
func BuildFeedMessage(entities []types.Entity) *gtfsrt.FeedMessage {
	return buildFeedMessage(entities)
}

// BuildPerDatasource groups entities by datasource and returns one FeedMessage per key.
func BuildPerDatasource(entities []types.Entity) map[string]*gtfsrt.FeedMessage {
	return buildPerDatasource(entities)
}
