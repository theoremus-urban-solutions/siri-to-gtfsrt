package siri

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
