package siri

// Minimal structs sufficient for ET/VM/SX mapping. More fields can be added as needed.

type ServiceDelivery struct {
	EstimatedTimetableDeliveries []EstimatedTimetableDelivery `xml:"EstimatedTimetableDelivery"`
	VehicleMonitoringDeliveries  []VehicleMonitoringDelivery  `xml:"VehicleMonitoringDelivery"`
	SituationExchangeDeliveries  []SituationExchangeDelivery  `xml:"SituationExchangeDelivery"`
}
