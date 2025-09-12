package siri

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
