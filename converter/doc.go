// Package converter provides SIRI to GTFS-Realtime conversion logic.
//
// This package contains the core business logic for transforming SIRI
// (Service Interface for Real Time Information) messages into GTFS-Realtime
// feed messages. It handles:
//   - Vehicle Monitoring → Vehicle Positions
//   - Estimated Timetable → Trip Updates
//   - Situation Exchange → Service Alerts
//
// The package is designed to be stateless and functional, accepting parsed
// SIRI data and returning GTFS-RT entities.
//
// Example:
//
//	entities, err := converter.ConvertSIRI(serviceDelivery, converter.DefaultOptions())
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	feedMsg := converter.BuildFeedMessage(entities)
//	// Or organize by datasource:
//	byDatasource := converter.BuildPerDatasource(entities)
package converter
