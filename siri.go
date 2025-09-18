package main

import (
	"time"
)

// ServiceDelivery and SIRI domain types

type ServiceDelivery struct {
	EstimatedTimetableDeliveries []EstimatedTimetableDelivery `xml:"EstimatedTimetableDelivery"`
	VehicleMonitoringDeliveries  []VehicleMonitoringDelivery  `xml:"VehicleMonitoringDelivery"`
	SituationExchangeDeliveries  []SituationExchangeDelivery  `xml:"SituationExchangeDelivery"`
}

// Vehicle Monitoring (VM)

type VehicleMonitoringDelivery struct {
	VehicleActivities []VehicleActivity `xml:"VehicleActivity"`
}

type VehicleActivity struct {
	RecordedAtTime          *string                  `xml:"RecordedAtTime"`
	ValidUntilTime          *string                  `xml:"ValidUntilTime"`
	ProgressBetweenStops    *ProgressBetweenStops    `xml:"ProgressBetweenStops"`
	MonitoredVehicleJourney *MonitoredVehicleJourney `xml:"MonitoredVehicleJourney"`
}

type ProgressBetweenStops struct {
	LinkDistance *float64 `xml:"LinkDistance"`
	Percentage   *float64 `xml:"Percentage"`
}

type MonitoredVehicleJourney struct {
	LineRef                  *string                  `xml:"LineRef"`
	VehicleRef               *string                  `xml:"VehicleRef"`
	DataSource               *string                  `xml:"DataSource"`
	VehicleLocation          *Location                `xml:"VehicleLocation"`
	Bearing                  *float32                 `xml:"Bearing"`
	Velocity                 *float64                 `xml:"Velocity"`
	MonitoredCall            *MonitoredCall           `xml:"MonitoredCall"`
	FramedVehicleJourneyRef  *FramedVehicleJourneyRef `xml:"FramedVehicleJourneyRef"`
	OriginAimedDepartureTime *string                  `xml:"OriginAimedDepartureTime"`
}

type Location struct {
	Longitude float64 `xml:"Longitude"`
	Latitude  float64 `xml:"Latitude"`
}

type MonitoredCall struct {
	StopPointRef          *string   `xml:"StopPointRef"`
	VehicleAtStop         *bool     `xml:"VehicleAtStop"`
	VehicleLocationAtStop *Location `xml:"VehicleLocationAtStop"`
	Order                 *int32    `xml:"Order"`
}

// Estimated Timetable (ET)

type EstimatedTimetableDelivery struct {
	EstimatedJourneyVersionFrames []EstimatedJourneyVersionFrame `xml:"EstimatedJourneyVersionFrame"`
}

type EstimatedJourneyVersionFrame struct {
	EstimatedVehicleJourneys []EstimatedVehicleJourney `xml:"EstimatedVehicleJourney"`
}

type EstimatedVehicleJourney struct {
	RecordedAtTime           *string                  `xml:"RecordedAtTime"`
	LineRef                  *string                  `xml:"LineRef"`
	FramedVehicleJourneyRef  *FramedVehicleJourneyRef `xml:"FramedVehicleJourneyRef"`
	DatedVehicleJourneyRef   *string                  `xml:"DatedVehicleJourneyRef"`
	VehicleRef               *string                  `xml:"VehicleRef"`
	OriginAimedDepartureTime *string                  `xml:"OriginAimedDepartureTime"`
	DataSource               *string                  `xml:"DataSource"`
	RecordedCalls            []RecordedCall           `xml:"RecordedCalls>RecordedCall"`
	EstimatedCalls           []EstimatedCall          `xml:"EstimatedCalls>EstimatedCall"`
}

type FramedVehicleJourneyRef struct {
	DataFrameRef           *string `xml:"DataFrameRef"`
	DatedVehicleJourneyRef *string `xml:"DatedVehicleJourneyRef"`
}

type RecordedCall struct {
	StopPointRef          *string `xml:"StopPointRef"`
	Order                 *int32  `xml:"Order"`
	AimedArrivalTime      *string `xml:"AimedArrivalTime"`
	ExpectedArrivalTime   *string `xml:"ExpectedArrivalTime"`
	ActualArrivalTime     *string `xml:"ActualArrivalTime"`
	AimedDepartureTime    *string `xml:"AimedDepartureTime"`
	ExpectedDepartureTime *string `xml:"ExpectedDepartureTime"`
	ActualDepartureTime   *string `xml:"ActualDepartureTime"`
}

type EstimatedCall struct {
	StopPointRef          *string `xml:"StopPointRef"`
	Order                 *int32  `xml:"Order"`
	AimedArrivalTime      *string `xml:"AimedArrivalTime"`
	ExpectedArrivalTime   *string `xml:"ExpectedArrivalTime"`
	AimedDepartureTime    *string `xml:"AimedDepartureTime"`
	ExpectedDepartureTime *string `xml:"ExpectedDepartureTime"`
}

// Situation Exchange (SX)

type SituationExchangeDelivery struct {
	Situations []PtSituationElement `xml:"Situations>PtSituationElement"`
}

type PtSituationElement struct {
	ParticipantRef  *string          `xml:"ParticipantRef"`
	SituationNumber *string          `xml:"SituationNumber"`
	ValidityPeriods []ValidityPeriod `xml:"ValidityPeriod"`
	Summaries       []TranslatedText `xml:"Summary"`
	Descriptions    []TranslatedText `xml:"Description"`
	Affects         *Affects         `xml:"Affects"`
	InfoLinks       []InfoLink       `xml:"InfoLinks>InfoLink"`
}

type ValidityPeriod struct {
	StartTime *string `xml:"StartTime"`
	EndTime   *string `xml:"EndTime"`
}

type TranslatedText struct {
	Value string `xml:",chardata"`
	Lang  string `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
}

type InfoLink struct {
	Uri string `xml:"Uri"`
}

type Affects struct {
	StopPoints      []AffectedStopPoint      `xml:"StopPoints>AffectedStopPoint"`
	VehicleJourneys []AffectedVehicleJourney `xml:"VehicleJourneys>AffectedVehicleJourney"`
	Networks        []AffectedNetwork        `xml:"Networks>AffectedNetwork"`
	StopPlaces      []AffectedStopPlace      `xml:"StopPlaces>AffectedStopPlace"`
}

type AffectedStopPoint struct {
	StopPointRef *string `xml:"StopPointRef"`
}

type AffectedVehicleJourney struct {
	LineRef                  *string                  `xml:"LineRef"`
	VehicleJourneyRefs       []string                 `xml:"VehicleJourneyRef"`
	DatedVehicleJourneyRefs  []string                 `xml:"DatedVehicleJourneyRef"`
	FramedVehicleJourneyRef  *FramedVehicleJourneyRef `xml:"FramedVehicleJourneyRef"`
	Routes                   []AffectedRoute          `xml:"Routes>AffectedRoute"`
	OriginAimedDepartureTime *string                  `xml:"OriginAimedDepartureTime"`
}

type AffectedRoute struct {
	StopPoints StopPoints `xml:"StopPoints"`
}

type StopPoints struct {
	StopPoints []AffectedStopPoint `xml:"AffectedStopPoint"`
}

type AffectedNetwork struct {
	AffectedLines []AffectedLine `xml:"AffectedLine"`
}

type AffectedLine struct {
	LineRef *string         `xml:"LineRef"`
	Routes  []AffectedRoute `xml:"Routes>AffectedRoute"`
}

type AffectedStopPlace struct {
	StopPlaceRef *string `xml:"StopPlaceRef"`
}

// Time helpers

// ParseISOTime attempts to parse common ISO-8601/RFC3339 timestamps with offsets.
func ParseISOTime(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05Z07:00",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// Latest returns the latest of t1 and t2, handling zero values.
func Latest(t1, t2 time.Time) time.Time {
	if t1.IsZero() {
		return t2
	}
	if t2.IsZero() {
		return t1
	}
	if t2.After(t1) {
		return t2
	}
	return t1
}

// FormatDateYYYYMMDD formats a time into YYYYMMDD string.
func FormatDateYYYYMMDD(t time.Time) string { return t.Format("20060102") }
