package formatter_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/theoremus-urban-solutions/siri-to-gtfsrt/formatter"
)

// Example demonstrates basic SIRI XML decoding
func Example() {
	xmlData := `<Siri version="2.0" xmlns="http://www.siri.org.uk/siri">
  <ServiceDelivery>
    <VehicleMonitoringDelivery>
      <VehicleActivity>
        <MonitoredVehicleJourney>
          <LineRef>TEST:Line:1</LineRef>
        </MonitoredVehicleJourney>
      </VehicleActivity>
    </VehicleMonitoringDelivery>
  </ServiceDelivery>
</Siri>`

	reader := bytes.NewReader([]byte(xmlData))
	sd, err := formatter.DecodeSIRI(reader)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("VehicleMonitoringDeliveries: %d\n", len(sd.VehicleMonitoringDeliveries))
	// Output: VehicleMonitoringDeliveries: 1
}

// ExampleDecodeSIRIFromBase64 demonstrates streaming base64 decoding
func ExampleDecodeSIRIFromBase64() {
	xmlData := `<Siri version="2.0" xmlns="http://www.siri.org.uk/siri">
  <ServiceDelivery>
    <VehicleMonitoringDelivery>
      <VehicleActivity>
        <MonitoredVehicleJourney>
          <LineRef>TEST:Line:1</LineRef>
        </MonitoredVehicleJourney>
      </VehicleActivity>
    </VehicleMonitoringDelivery>
  </ServiceDelivery>
</Siri>`

	// Encode to base64
	base64Data := base64.StdEncoding.EncodeToString([]byte(xmlData))

	// Stream decode: base64 → XML parser (no intermediate buffering)
	reader := bytes.NewReader([]byte(base64Data))
	sd, err := formatter.DecodeSIRIFromBase64(reader)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("VehicleMonitoringDeliveries: %d\n", len(sd.VehicleMonitoringDeliveries))
	// Output: VehicleMonitoringDeliveries: 1
}
