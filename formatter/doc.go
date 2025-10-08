// Package formatter provides input/output formatting for SIRI messages.
//
// This package handles serialization and deserialization of SIRI data
// in various formats. Currently supports:
//   - XML: SIRI XML message decoding
//
// Future support may include:
//   - JSON: SIRI JSON encoding/decoding
//   - Custom formatters for specific data sources
//
// Example:
//
//	file, _ := os.Open("siri-vm.xml")
//	defer file.Close()
//
//	serviceDelivery, err := formatter.DecodeSIRI(file)
//	if err != nil {
//	    log.Fatal(err)
//	}
package formatter
