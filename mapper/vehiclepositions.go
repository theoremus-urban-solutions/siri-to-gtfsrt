package mapper

import (
	"golang/gtfsrt"
	"golang/siri"
	"golang/types"
	"time"
)

// MapVMToVehiclePosition converts a SIRI VehicleActivity into a GTFS-RT entity.
func MapVMToVehiclePosition(va *siri.VehicleActivity, opts types.Options) *types.Entity {
	if va == nil || va.MonitoredVehicleJourney == nil {
		return nil
	}
	mvj := va.MonitoredVehicleJourney
	if mvj.VehicleLocation == nil || (mvj.FramedVehicleJourneyRef == nil && mvj.VehicleRef == nil) {
		return nil
	}

	// Build ID: prefer vehicleRef; else tripId-startDate
	var id string
	if mvj.VehicleRef != nil && *mvj.VehicleRef != "" {
		id = *mvj.VehicleRef
	} else if mvj.FramedVehicleJourneyRef != nil && mvj.FramedVehicleJourneyRef.DatedVehicleJourneyRef != nil {
		id = *mvj.FramedVehicleJourneyRef.DatedVehicleJourneyRef
		if mvj.OriginAimedDepartureTime != nil {
			if t, ok := siri.ParseISOTime(*mvj.OriginAimedDepartureTime); ok {
				id = id + "-" + siri.FormatDateYYYYMMDD(t)
			}
		}
	}
	if id == "" {
		return nil
	}

	// TTL from ValidUntilTime or grace
	ttl := opts.VMGracePeriod
	if va.ValidUntilTime != nil {
		if t, ok := siri.ParseISOTime(*va.ValidUntilTime); ok {
			d := time.Until(t)
			if d > 0 {
				ttl = d
			}
		}
	}

	// Build minimal VehiclePosition for JSON output
	ent := &gtfsrt.FeedEntity{Id: &id}
	vp := &gtfsrt.VehiclePosition{}
	if mvj.FramedVehicleJourneyRef != nil && mvj.FramedVehicleJourneyRef.DatedVehicleJourneyRef != nil {
		td := &gtfsrt.TripDescriptor{TripId: *mvj.FramedVehicleJourneyRef.DatedVehicleJourneyRef}
		if mvj.OriginAimedDepartureTime != nil {
			if t, ok := siri.ParseISOTime(*mvj.OriginAimedDepartureTime); ok {
				td.StartDate = siri.FormatDateYYYYMMDD(t)
			}
		}
		if mvj.LineRef != nil {
			td.RouteId = *mvj.LineRef
		}
		vp.Trip = td
	}
	if mvj.VehicleRef != nil && *mvj.VehicleRef != "" {
		vp.Vehicle = &gtfsrt.VehicleDescriptor{Id: *mvj.VehicleRef}
	}
	if mvj.VehicleLocation != nil {
		pos := &gtfsrt.Position{Latitude: float32(mvj.VehicleLocation.Latitude), Longitude: float32(mvj.VehicleLocation.Longitude)}
		if mvj.Bearing != nil {
			pos.Bearing = mvj.Bearing
		}
		if mvj.Velocity != nil {
			sp := float32(*mvj.Velocity)
			pos.Speed = &sp
		}
		vp.Position = pos
	}
	if va.RecordedAtTime != nil {
		if t, ok := siri.ParseISOTime(*va.RecordedAtTime); ok {
			ts := t.Unix()
			vp.Timestamp = &ts
		}
	}
	ent.Vehicle = vp

	return &types.Entity{ID: id, Datasource: deref(mvj.DataSource), Message: ent, TTL: ttl}
}
