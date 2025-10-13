package formatter_test

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/theoremus-urban-solutions/siri-to-gtfsrt/formatter"
)

func TestDecodeSIRIFromBase64(t *testing.T) {
	// Minimal SIRI XML for testing
	xmlData := `<Siri version="2.0" xmlns="http://www.siri.org.uk/siri">
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

	// Encode to base64
	base64Data := base64.StdEncoding.EncodeToString([]byte(xmlData))

	t.Run("decode valid base64 SIRI XML", func(t *testing.T) {
		reader := bytes.NewReader([]byte(base64Data))
		sd, err := formatter.DecodeSIRIFromBase64(reader)

		if err != nil {
			t.Fatalf("DecodeSIRIFromBase64 failed: %v", err)
		}

		if sd == nil {
			t.Fatal("ServiceDelivery is nil")
		}

		if len(sd.VehicleMonitoringDeliveries) != 1 {
			t.Errorf("Expected 1 VehicleMonitoringDelivery, got %d", len(sd.VehicleMonitoringDeliveries))
		}

		if len(sd.VehicleMonitoringDeliveries) > 0 {
			vmd := sd.VehicleMonitoringDeliveries[0]
			if len(vmd.VehicleActivities) != 1 {
				t.Errorf("Expected 1 VehicleActivity, got %d", len(vmd.VehicleActivities))
			}

			if len(vmd.VehicleActivities) > 0 {
				va := vmd.VehicleActivities[0]
				if va.MonitoredVehicleJourney == nil || va.MonitoredVehicleJourney.LineRef == nil {
					t.Error("MonitoredVehicleJourney or LineRef is nil")
				} else if *va.MonitoredVehicleJourney.LineRef != "TEST:Line:1" {
					t.Errorf("Expected LineRef 'TEST:Line:1', got '%s'", *va.MonitoredVehicleJourney.LineRef)
				}
			}
		}
	})

	t.Run("decode empty base64", func(t *testing.T) {
		reader := bytes.NewReader([]byte(""))
		sd, err := formatter.DecodeSIRIFromBase64(reader)

		// Should not error, but return empty ServiceDelivery
		if err != nil {
			t.Fatalf("DecodeSIRIFromBase64 failed on empty input: %v", err)
		}

		if sd == nil {
			t.Fatal("ServiceDelivery is nil")
		}
	})

	t.Run("streaming behavior - large input", func(t *testing.T) {
		// Create a larger XML by repeating vehicle activities
		var sb strings.Builder
		sb.WriteString(`<Siri version="2.0" xmlns="http://www.siri.org.uk/siri">
  <ServiceDelivery>
    <VehicleMonitoringDelivery>`)

		// Add 100 vehicle activities
		for i := 0; i < 100; i++ {
			sb.WriteString(`
      <VehicleActivity>
        <RecordedAtTime>2025-09-12T10:00:00+00:00</RecordedAtTime>
        <MonitoredVehicleJourney>
          <LineRef>TEST:Line:1</LineRef>
        </MonitoredVehicleJourney>
      </VehicleActivity>`)
		}

		sb.WriteString(`
    </VehicleMonitoringDelivery>
  </ServiceDelivery>
</Siri>`)

		largeXML := sb.String()
		base64Large := base64.StdEncoding.EncodeToString([]byte(largeXML))

		reader := bytes.NewReader([]byte(base64Large))
		sd, err := formatter.DecodeSIRIFromBase64(reader)

		if err != nil {
			t.Fatalf("DecodeSIRIFromBase64 failed on large input: %v", err)
		}

		if sd == nil {
			t.Fatal("ServiceDelivery is nil")
		}

		if len(sd.VehicleMonitoringDeliveries) != 1 {
			t.Errorf("Expected 1 VehicleMonitoringDelivery, got %d", len(sd.VehicleMonitoringDeliveries))
		}

		if len(sd.VehicleMonitoringDeliveries) > 0 {
			vmd := sd.VehicleMonitoringDeliveries[0]
			if len(vmd.VehicleActivities) != 100 {
				t.Errorf("Expected 100 VehicleActivities, got %d", len(vmd.VehicleActivities))
			}
		}
	})
}

func TestDecodeSIRIFromBase64_CompareWithDirect(t *testing.T) {
	// Test that streaming base64 decode produces same result as decode-then-parse
	xmlData := `<Siri version="2.0" xmlns="http://www.siri.org.uk/siri">
  <ServiceDelivery>
    <EstimatedTimetableDelivery>
      <EstimatedJourneyVersionFrame>
        <EstimatedVehicleJourney>
          <LineRef>TEST:Line:2</LineRef>
          <DatedVehicleJourneyRef>TEST:ServiceJourney:2</DatedVehicleJourneyRef>
        </EstimatedVehicleJourney>
      </EstimatedJourneyVersionFrame>
    </EstimatedTimetableDelivery>
  </ServiceDelivery>
</Siri>`

	base64Data := base64.StdEncoding.EncodeToString([]byte(xmlData))

	// Method 1: Stream base64 decode
	reader1 := bytes.NewReader([]byte(base64Data))
	sd1, err1 := formatter.DecodeSIRIFromBase64(reader1)
	if err1 != nil {
		t.Fatalf("DecodeSIRIFromBase64 failed: %v", err1)
	}

	// Method 2: Decode base64 first, then parse XML
	xmlBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		t.Fatalf("base64 decode failed: %v", err)
	}
	reader2 := bytes.NewReader(xmlBytes)
	sd2, err2 := formatter.DecodeSIRI(reader2)
	if err2 != nil {
		t.Fatalf("DecodeSIRI failed: %v", err2)
	}

	// Compare results
	if len(sd1.EstimatedTimetableDeliveries) != len(sd2.EstimatedTimetableDeliveries) {
		t.Errorf("EstimatedTimetableDeliveries count mismatch: %d vs %d",
			len(sd1.EstimatedTimetableDeliveries), len(sd2.EstimatedTimetableDeliveries))
	}

	if len(sd1.EstimatedTimetableDeliveries) > 0 && len(sd2.EstimatedTimetableDeliveries) > 0 {
		etd1 := sd1.EstimatedTimetableDeliveries[0]
		etd2 := sd2.EstimatedTimetableDeliveries[0]

		if len(etd1.EstimatedJourneyVersionFrames) != len(etd2.EstimatedJourneyVersionFrames) {
			t.Errorf("EstimatedJourneyVersionFrames count mismatch: %d vs %d",
				len(etd1.EstimatedJourneyVersionFrames), len(etd2.EstimatedJourneyVersionFrames))
		}

		if len(etd1.EstimatedJourneyVersionFrames) > 0 && len(etd2.EstimatedJourneyVersionFrames) > 0 {
			frame1 := etd1.EstimatedJourneyVersionFrames[0]
			frame2 := etd2.EstimatedJourneyVersionFrames[0]

			if len(frame1.EstimatedVehicleJourneys) != len(frame2.EstimatedVehicleJourneys) {
				t.Errorf("EstimatedVehicleJourneys count mismatch: %d vs %d",
					len(frame1.EstimatedVehicleJourneys), len(frame2.EstimatedVehicleJourneys))
			}

			if len(frame1.EstimatedVehicleJourneys) > 0 && len(frame2.EstimatedVehicleJourneys) > 0 {
				evj1 := frame1.EstimatedVehicleJourneys[0]
				evj2 := frame2.EstimatedVehicleJourneys[0]

				if (evj1.LineRef == nil) != (evj2.LineRef == nil) {
					t.Error("LineRef nil mismatch")
				} else if evj1.LineRef != nil && evj2.LineRef != nil && *evj1.LineRef != *evj2.LineRef {
					t.Errorf("LineRef mismatch: '%s' vs '%s'", *evj1.LineRef, *evj2.LineRef)
				}

				if (evj1.DatedVehicleJourneyRef == nil) != (evj2.DatedVehicleJourneyRef == nil) {
					t.Error("DatedVehicleJourneyRef nil mismatch")
				} else if evj1.DatedVehicleJourneyRef != nil && evj2.DatedVehicleJourneyRef != nil &&
					*evj1.DatedVehicleJourneyRef != *evj2.DatedVehicleJourneyRef {
					t.Errorf("DatedVehicleJourneyRef mismatch: '%s' vs '%s'",
						*evj1.DatedVehicleJourneyRef, *evj2.DatedVehicleJourneyRef)
				}
			}
		}
	}
}
