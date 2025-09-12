package mapper

import (
	"golang/gtfsrt"
	"golang/siri"
	"golang/types"
	"time"
)

// MapETToTripUpdate converts a SIRI EstimatedVehicleJourney into a GTFS-RT entity.
func MapETToTripUpdate(evj *siri.EstimatedVehicleJourney, opts types.Options) *types.Entity {
	if evj == nil || evj.FramedVehicleJourneyRef == nil || evj.FramedVehicleJourneyRef.DatedVehicleJourneyRef == nil {
		return nil
	}

	// Build minimal entity id: tripId + startDate if available
	tripId := *evj.FramedVehicleJourneyRef.DatedVehicleJourneyRef
	var startDate string
	if evj.OriginAimedDepartureTime != nil {
		if t, ok := siri.ParseISOTime(*evj.OriginAimedDepartureTime); ok {
			startDate = siri.FormatDateYYYYMMDD(t)
		}
	}

	// Build entity id and message
	id := tripId
	if startDate != "" {
		id = tripId + "-" + startDate
	}

	// TTL: max across recorded/estimated calls
	var latest time.Time
	for _, rc := range evj.RecordedCalls {
		times := []*string{rc.ActualArrivalTime, rc.ExpectedArrivalTime, rc.AimedArrivalTime, rc.ActualDepartureTime, rc.ExpectedDepartureTime, rc.AimedDepartureTime}
		for _, ts := range times {
			if ts != nil {
				if t, ok := siri.ParseISOTime(*ts); ok {
					latest = siri.Latest(latest, t)
				}
			}
		}
	}
	for _, ec := range evj.EstimatedCalls {
		times := []*string{ec.ExpectedArrivalTime, ec.AimedArrivalTime, ec.ExpectedDepartureTime, ec.AimedDepartureTime}
		for _, ts := range times {
			if ts != nil {
				if t, ok := siri.ParseISOTime(*ts); ok {
					latest = siri.Latest(latest, t)
				}
			}
		}
	}
	ttl := opts.VMGracePeriod
	if !latest.IsZero() {
		d := time.Until(latest)
		if d > 0 {
			ttl = d
		}
	}

	// Build TripUpdate message
	ent := &gtfsrt.FeedEntity{Id: &id}
	tu := &gtfsrt.TripUpdate{}
	td := &gtfsrt.TripDescriptor{TripId: tripId}
	if evj.LineRef != nil {
		td.RouteId = *evj.LineRef
	}
	if evj.OriginAimedDepartureTime != nil {
		if t, ok := siri.ParseISOTime(*evj.OriginAimedDepartureTime); ok {
			td.StartDate = siri.FormatDateYYYYMMDD(t)
			td.StartTime = t.Format("15:04:05")
		}
	}
	tu.Trip = td
	if evj.VehicleRef != nil && *evj.VehicleRef != "" {
		tu.Vehicle = &gtfsrt.VehicleDescriptor{Id: *evj.VehicleRef}
	}

	stopSeq := int32(0)
	for _, rc := range evj.RecordedCalls {
		stu := gtfsrt.StopTimeUpdate{}
		if rc.StopPointRef != nil {
			stu.StopId = *rc.StopPointRef
		}
		if rc.Order != nil && *rc.Order > 0 {
			stu.StopSequence = *rc.Order - 1
		} else {
			stu.StopSequence = stopSeq
		}
		if rc.AimedArrivalTime != nil {
			if d := delaySeconds(rc.AimedArrivalTime, rc.ActualArrivalTime, rc.ExpectedArrivalTime); d != nil {
				stu.Arrival = &gtfsrt.StopTimeEvent{Delay: d}
			}
		}
		if rc.AimedDepartureTime != nil {
			if d := delaySeconds(rc.AimedDepartureTime, rc.ActualDepartureTime, rc.ExpectedDepartureTime); d != nil {
				stu.Departure = &gtfsrt.StopTimeEvent{Delay: d}
			}
		}
		tu.StopTimeUpdate = append(tu.StopTimeUpdate, stu)
		stopSeq++
	}
	for _, ec := range evj.EstimatedCalls {
		stu := gtfsrt.StopTimeUpdate{}
		if ec.StopPointRef != nil {
			stu.StopId = *ec.StopPointRef
		}
		if ec.Order != nil && *ec.Order > 0 {
			stu.StopSequence = *ec.Order - 1
		} else {
			stu.StopSequence = stopSeq
		}
		if ec.AimedArrivalTime != nil && ec.ExpectedArrivalTime != nil {
			if d := diffSeconds(*ec.AimedArrivalTime, *ec.ExpectedArrivalTime); d != nil {
				stu.Arrival = &gtfsrt.StopTimeEvent{Delay: d}
			}
		}
		if ec.AimedDepartureTime != nil && ec.ExpectedDepartureTime != nil {
			if d := diffSeconds(*ec.AimedDepartureTime, *ec.ExpectedDepartureTime); d != nil {
				stu.Departure = &gtfsrt.StopTimeEvent{Delay: d}
			}
		}
		tu.StopTimeUpdate = append(tu.StopTimeUpdate, stu)
		stopSeq++
	}

	ent.TripUpdate = tu

	return &types.Entity{
		ID:         id,
		Datasource: deref(evj.DataSource),
		Message:    ent,
		TTL:        ttl,
	}
}

func deref[T ~string](p *T) string {
	if p == nil {
		return ""
	}
	return string(*p)
}

// helpers for delays
func delaySeconds(aimed, actual, expected *string) *int32 {
	if actual != nil {
		return diffSeconds(*aimed, *actual)
	}
	if expected != nil {
		return diffSeconds(*aimed, *expected)
	}
	return nil
}

func diffSeconds(aimed, updated string) *int32 {
	at, ok1 := siri.ParseISOTime(aimed)
	ut, ok2 := siri.ParseISOTime(updated)
	if !ok1 || !ok2 {
		return nil
	}
	d := int32(ut.Sub(at).Seconds())
	return &d
}
