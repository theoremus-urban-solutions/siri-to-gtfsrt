package gtfsrt

import (
	"fmt"
	"time"
)

// NOTE: Placeholder GTFS-RT-like types to allow compilation and JSON output
// before wiring real protobuf bindings. Replace with MobilityData bindings later.

type FeedHeader struct {
	GtfsRealtimeVersion *string `json:"gtfs_realtime_version,omitempty"`
	Incrementality      *int32  `json:"incrementality,omitempty"`
	Timestamp           string  `json:"timestamp,omitempty"`
}

const (
	IncrementalityFullDataset int32 = 1
)

type FeedMessage struct {
	Entity []*FeedEntity `json:"entity,omitempty"`
	Header *FeedHeader   `json:"header,omitempty"`
}

type FeedEntity struct {
	Alert      *Alert           `json:"alert,omitempty"`
	Id         *string          `json:"id,omitempty"`
	IsDeleted  *bool            `json:"is_deleted,omitempty"`
	TripUpdate *TripUpdate      `json:"trip_update,omitempty"`
	Vehicle    *VehiclePosition `json:"vehicle,omitempty"`
}

// TripUpdate

type TripUpdate struct {
	StopTimeUpdate []StopTimeUpdate   `json:"stop_time_update,omitempty"`
	Timestamp      string             `json:"timestamp,omitempty"`
	Trip           *TripDescriptor    `json:"trip,omitempty"`
	Vehicle        *VehicleDescriptor `json:"vehicle,omitempty"`
}

type TripDescriptor struct {
	RouteId              string `json:"route_id,omitempty"`
	ScheduleRelationship *int32 `json:"schedule_relationship,omitempty"`
	TripId               string `json:"trip_id,omitempty"`
	StartDate            string `json:"start_date,omitempty"`
	StartTime            string `json:"start_time,omitempty"`
}

type VehicleDescriptor struct {
	Id string `json:"id,omitempty"`
}

type StopTimeUpdate struct {
	Arrival              *StopTimeEvent `json:"arrival,omitempty"`
	Departure            *StopTimeEvent `json:"departure,omitempty"`
	ScheduleRelationship *int32         `json:"schedule_relationship,omitempty"`
	StopId               string         `json:"stop_id,omitempty"`
	StopSequence         int32          `json:"stop_sequence,omitempty"`
}

type StopTimeEvent struct {
	Time        string `json:"time,omitempty"`
	Uncertainty *int32 `json:"uncertainty,omitempty"`
	Delay       *int32 `json:"delay,omitempty"`
}

// VehiclePosition

type VehiclePosition struct {
	CongestionLevel     *string            `json:"congestion_level,omitempty"`
	CurrentStatus       string             `json:"current_status,omitempty"`
	CurrentStopSequence *int32             `json:"current_stop_sequence,omitempty"`
	OccupancyStatus     *string            `json:"occupancy_status,omitempty"`
	Position            *Position          `json:"position,omitempty"`
	StopId              *string            `json:"stop_id,omitempty"`
	Timestamp           *int64             `json:"timestamp,omitempty"`
	Trip                *TripDescriptor    `json:"trip,omitempty"`
	Vehicle             *VehicleDescriptor `json:"vehicle,omitempty"`
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
	ActivePeriod    []TimeRange       `json:"active_period,omitempty"`
	Cause           *int32            `json:"cause,omitempty"`
	DescriptionText *TranslatedString `json:"description_text,omitempty"`
	Effect          *int32            `json:"effect,omitempty"`
	HeaderText      *TranslatedString `json:"header_text,omitempty"`
	InformedEntity  []EntitySelector  `json:"informed_entity,omitempty"`
	Url             *TranslatedString `json:"url,omitempty"`
	// Raw SIRI fields for PBF mapping (not emitted in JSON)
	Severity *string `json:"-"`
}

type TranslatedString struct {
	Translation []Translation `json:"translation,omitempty"`
}

type Translation struct {
	Language *string `json:"language,omitempty"`
	Text     string  `json:"text"`
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
	ts := fmt.Sprintf("%d", time.Now().Unix())
	hdr := &FeedHeader{
		Timestamp:           ts,
		GtfsRealtimeVersion: stringPtr("2.0"),
		Incrementality:      &[]int32{0}[0], // FULL_DATASET = 0
	}
	return hdr
}

// NewFeedMessage returns a FeedMessage with a populated header.
func NewFeedMessage() *FeedMessage {
	return &FeedMessage{Header: NewFeedMessageHeader()}
}

func stringPtr(s string) *string { return &s }
