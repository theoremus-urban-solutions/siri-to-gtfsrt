package formatter

import (
	"encoding/base64"
	"io"

	"github.com/theoremus-urban-solutions/siri-to-gtfsrt/siri"
)

// DecodeSIRIFromBase64 reads base64-encoded SIRI XML and returns a populated ServiceDelivery.
// This function streams the base64 decoding directly into XML parsing, avoiding intermediate
// buffering of the full XML string in memory.
//
// Example usage:
//
//	// Server receives base64-encoded XML
//	base64Reader := bytes.NewReader(base64EncodedData)
//	sd, err := formatter.DecodeSIRIFromBase64(base64Reader)
func DecodeSIRIFromBase64(r io.Reader) (*siri.ServiceDelivery, error) {
	// Wrap the input reader with a streaming base64 decoder
	// This decodes base64 on-the-fly as the XML parser reads data
	base64Decoder := base64.NewDecoder(base64.StdEncoding, r)

	// Pass the streaming decoder directly to the XML parser
	return DecodeSIRI(base64Decoder)
}
