// Package siri provides SIRI (Service Interface for Real Time Information) domain types.
//
// SIRI is a European standard for real-time public transport information exchange.
// This package defines the Go structures for SIRI data including:
//   - ServiceDelivery: Root container for SIRI messages
//   - VehicleMonitoring: Real-time vehicle position updates
//   - EstimatedTimetable: Predicted arrival/departure times
//   - SituationExchange: Service alerts and disruptions
//
// All types in this package are pure domain models with no business logic.
// Time utilities for parsing SIRI timestamps are also provided.
//
// Example:
//
//	var sd siri.ServiceDelivery
//	// ... populate from XML ...
//	for _, delivery := range sd.VehicleMonitoringDeliveries {
//	    for _, activity := range delivery.VehicleActivities {
//	        // Process vehicle activity
//	    }
//	}
package siri
