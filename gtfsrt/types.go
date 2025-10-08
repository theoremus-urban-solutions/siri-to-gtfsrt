package gtfsrt

import "time"

// NOTE: Placeholder GTFS-RT-like types to allow compilation and JSON output
// before wiring real protobuf bindings. Replace with MobilityData bindings later.

type FeedHeader struct {
	Timestamp           *uint64 `json:"timestamp,omitempty"`
	GtfsRealtimeVersion *string `json:"gtfs_realtime_version,omitempty"`
	Incrementality      *int32  `json:"incrementality,omitempty"`
}

const (
	IncrementalityFullDataset int32 = 1
)

type FeedMessage struct {
	Header *FeedHeader   `json:"header,omitempty"`
	Entity []*FeedEntity `json:"entity,omitempty"`
}

type FeedEntity struct {
	Id         *string          `json:"id,omitempty"`
	TripUpdate *TripUpdate      `json:"trip_update,omitempty"`
	Vehicle    *VehiclePosition `json:"vehicle,omitempty"`
	Alert      *Alert           `json:"alert,omitempty"`
}

// TripUpdate

type TripUpdate struct {
	Trip           *TripDescriptor    `json:"trip,omitempty"`
	Vehicle        *VehicleDescriptor `json:"vehicle,omitempty"`
	StopTimeUpdate []StopTimeUpdate   `json:"stop_time_update,omitempty"`
}

type TripDescriptor struct {
	TripId    string `json:"trip_id,omitempty"`
	RouteId   string `json:"route_id,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	StartTime string `json:"start_time,omitempty"`
}

type VehicleDescriptor struct {
	Id string `json:"id,omitempty"`
}

type StopTimeUpdate struct {
	StopId       string         `json:"stop_id,omitempty"`
	StopSequence int32          `json:"stop_sequence,omitempty"`
	Arrival      *StopTimeEvent `json:"arrival,omitempty"`
	Departure    *StopTimeEvent `json:"departure,omitempty"`
}

type StopTimeEvent struct {
	Delay *int32 `json:"delay,omitempty"`
}

// VehiclePosition

type VehiclePosition struct {
	Trip                *TripDescriptor    `json:"trip,omitempty"`
	Vehicle             *VehicleDescriptor `json:"vehicle,omitempty"`
	Position            *Position          `json:"position,omitempty"`
	Timestamp           *int64             `json:"timestamp,omitempty"`
	CurrentStatus       string             `json:"current_status,omitempty"`
	StopId              *string            `json:"stop_id,omitempty"`
	CurrentStopSequence *int32             `json:"current_stop_sequence,omitempty"`
	OccupancyStatus     *string            `json:"occupancy_status,omitempty"`
	CongestionLevel     *string            `json:"congestion_level,omitempty"`
}

type Position struct {
	Latitude  float32  `json:"latitude,omitempty"`
	Longitude float32  `json:"longitude,omitempty"`
	Bearing   *float32 `json:"bearing,omitempty"`
	Speed     *float32 `json:"speed,omitempty"`
	Odometer  *float64 `json:"odometer,omitempty"`
}

// Alerts

type Alert struct {
	HeaderText      *TranslatedString `json:"header_text,omitempty"`
	DescriptionText *TranslatedString `json:"description_text,omitempty"`
	ActivePeriod    []TimeRange       `json:"active_period,omitempty"`
	InformedEntity  []EntitySelector  `json:"informed_entity,omitempty"`
	Url             *TranslatedString `json:"url,omitempty"`
	// Raw SIRI fields for PBF mapping (not emitted in JSON)
	Cause    *string `json:"-"`
	Effect   *string `json:"-"`
	Severity *string `json:"-"`
}

type TranslatedString struct {
	Translation []Translation `json:"translation,omitempty"`
}

type Translation struct {
	Text     string  `json:"text,omitempty"`
	Language *string `json:"language,omitempty"`
}

type TimeRange struct {
	Start *int64 `json:"start,omitempty"`
	End   *int64 `json:"end,omitempty"`
}

type EntitySelector struct {
	RouteId *string         `json:"route_id,omitempty"`
	StopId  *string         `json:"stop_id,omitempty"`
	Trip    *TripDescriptor `json:"trip,omitempty"`
}

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
