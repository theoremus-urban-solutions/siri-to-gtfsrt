// Package gtfsrt provides GTFS-Realtime type definitions and protobuf operations.
//
// GTFS-Realtime is a standard for providing real-time transit updates.
// This package provides:
//   - Internal type definitions that mirror the protobuf schema
//   - Conversion functions to official MobilityData protobuf types
//   - Serialization to Protocol Buffer format
//   - Utilities for unmarshaling and JSON conversion
//
// The package uses internal types during conversion before final serialization
// to allow for flexible manipulation before protobuf encoding.
//
// Example:
//
//	msg := &gtfsrt.FeedMessage{
//	    Header: gtfsrt.NewFeedMessageHeader(),
//	    Entity: []*gtfsrt.FeedEntity{...},
//	}
//	pbfBytes, _ := gtfsrt.MarshalPBF(msg)
package gtfsrt
