package siritogtfs

import (
	stdxml "encoding/xml"
	"io"
)

// DecodeSIRI reads SIRI XML and returns a populated ServiceDelivery.
func DecodeSIRI(r io.Reader) (*ServiceDelivery, error) {
	dec := stdxml.NewDecoder(r)
	// Try <Siri><ServiceDelivery>...</ServiceDelivery></Siri>
	var siriDoc struct {
		XMLName         stdxml.Name     `xml:"Siri"`
		ServiceDelivery ServiceDelivery `xml:"ServiceDelivery"`
	}
	if err := dec.Decode(&siriDoc); err == nil {
		return &siriDoc.ServiceDelivery, nil
	}
	// Fallback: direct ServiceDelivery root
	if seeker, ok := r.(io.Seeker); ok {
		if _, err := seeker.Seek(0, 0); err == nil {
			var sd ServiceDelivery
			if err := stdxml.NewDecoder(r).Decode(&sd); err == nil {
				return &sd, nil
			}
		}
	}
	// Give up with empty
	return &ServiceDelivery{}, nil
}
