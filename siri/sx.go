package siri

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
