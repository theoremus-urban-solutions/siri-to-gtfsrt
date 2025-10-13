// Package formatter provides input/output formatting for SIRI messages.
//
// This package handles serialization and deserialization of SIRI data
// in various formats. Currently supports:
//   - XML: SIRI XML message decoding
//   - Base64-encoded XML: Streaming decoder for optimal performance
//
// Future support may include:
//   - JSON: SIRI JSON encoding/decoding
//   - Custom formatters for specific data sources
//
// # Basic XML Decoding
//
//	file, _ := os.Open("siri-vm.xml")
//	defer file.Close()
//
//	serviceDelivery, err := formatter.DecodeSIRI(file)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Streaming Base64 Decoding
//
// For optimal performance when receiving base64-encoded XML:
//
//	// Streams base64 → XML parser without intermediate buffering
//	serviceDelivery, err := formatter.DecodeSIRIFromBase64(base64Reader)
//	if err != nil {
//	    log.Fatal(err)
//	}
package formatter
