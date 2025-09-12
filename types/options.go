package types

import "time"

// Options controls conversion behavior and filters.
type Options struct {
	ETWhitelist []string
	VMWhitelist []string
	SXWhitelist []string

	CloseToNextStopPercentage int
	CloseToNextStopDistance   int

	VMGracePeriod time.Duration
}

func DefaultOptions() Options {
	return Options{
		CloseToNextStopPercentage: 95,
		CloseToNextStopDistance:   500,
		VMGracePeriod:             5 * time.Minute,
	}
}
