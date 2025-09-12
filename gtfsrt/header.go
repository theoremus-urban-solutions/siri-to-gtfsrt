package gtfsrt

import "time"

// NewFeedMessageHeader creates a GTFS-RT header matching Java defaults.
func NewFeedMessageHeader() *FeedHeader {
	ts := uint64(time.Now().Unix())
	hdr := &FeedHeader{
		Timestamp:           &ts,
		GtfsRealtimeVersion: stringPtr("1.0"),
		Incrementality:      &[]int32{IncrementalityFullDataset}[0],
	}
	return hdr
}

// NewFeedMessage returns a FeedMessage with a populated header.
func NewFeedMessage() *FeedMessage {
	return &FeedMessage{Header: NewFeedMessageHeader()}
}

func stringPtr(s string) *string { return &s }
