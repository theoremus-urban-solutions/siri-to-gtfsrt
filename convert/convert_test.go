package convert_test

import (
	"strings"
	"testing"

	"golang/convert"
	sirixml "golang/siri/xml"
	"golang/types"
)

const minimalVM = `
<Siri version="2.0" xmlns="http://www.siri.org.uk/siri">
  <ServiceDelivery>
    <VehicleMonitoringDelivery>
      <VehicleActivity>
        <RecordedAtTime>2025-09-12T10:00:00+00:00</RecordedAtTime>
        <ValidUntilTime>2025-09-12T10:05:00+00:00</ValidUntilTime>
        <MonitoredVehicleJourney>
          <LineRef>TEST:Line:1</LineRef>
          <FramedVehicleJourneyRef>
            <DataFrameRef>2025-09-12</DataFrameRef>
            <DatedVehicleJourneyRef>TEST:ServiceJourney:1</DatedVehicleJourneyRef>
          </FramedVehicleJourneyRef>
          <OriginAimedDepartureTime>2025-09-12T09:55:00+00:00</OriginAimedDepartureTime>
          <VehicleLocation>
            <Longitude>10.0</Longitude>
            <Latitude>59.0</Latitude>
          </VehicleLocation>
        </MonitoredVehicleJourney>
      </VehicleActivity>
    </VehicleMonitoringDelivery>
  </ServiceDelivery>
</Siri>`

func TestConvertMinimalVM(t *testing.T) {
	sd, err := sirixml.Decode(strings.NewReader(minimalVM))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	entities, err := convert.ConvertSIRI(sd, types.DefaultOptions())
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if len(entities) == 0 {
		t.Fatalf("expected at least one entity")
	}
}
