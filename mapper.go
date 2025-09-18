package main

import "time"

// Mapping functions

// VM -> VehiclePosition
func MapVMToVehiclePosition(va *VehicleActivity, opts Options) *Entity {
	if va == nil || va.MonitoredVehicleJourney == nil {
		return nil
	}
	mvj := va.MonitoredVehicleJourney
	if mvj.VehicleLocation == nil || (mvj.FramedVehicleJourneyRef == nil && mvj.VehicleRef == nil) {
		return nil
	}

	var id string
	if mvj.VehicleRef != nil && *mvj.VehicleRef != "" {
		id = *mvj.VehicleRef
	} else if mvj.FramedVehicleJourneyRef != nil && mvj.FramedVehicleJourneyRef.DatedVehicleJourneyRef != nil {
		id = *mvj.FramedVehicleJourneyRef.DatedVehicleJourneyRef
		if mvj.OriginAimedDepartureTime != nil {
			if t, ok := ParseISOTime(*mvj.OriginAimedDepartureTime); ok {
				id = id + "-" + FormatDateYYYYMMDD(t)
			}
		}
	}
	if id == "" {
		return nil
	}

	ttl := opts.VMGracePeriod
	if va.ValidUntilTime != nil {
		if t, ok := ParseISOTime(*va.ValidUntilTime); ok {
			d := time.Until(t)
			if d > 0 {
				ttl = d
			}
		}
	}

	ent := &FeedEntity{Id: &id}
	vp := &VehiclePosition{}
	if mvj.FramedVehicleJourneyRef != nil && mvj.FramedVehicleJourneyRef.DatedVehicleJourneyRef != nil {
		td := &TripDescriptor{TripId: *mvj.FramedVehicleJourneyRef.DatedVehicleJourneyRef}
		if mvj.OriginAimedDepartureTime != nil {
			if t, ok := ParseISOTime(*mvj.OriginAimedDepartureTime); ok {
				td.StartDate = FormatDateYYYYMMDD(t)
			}
		}
		if mvj.LineRef != nil {
			td.RouteId = *mvj.LineRef
		}
		vp.Trip = td
	}
	if mvj.VehicleRef != nil && *mvj.VehicleRef != "" {
		vp.Vehicle = &VehicleDescriptor{Id: *mvj.VehicleRef}
	}
	if mvj.VehicleLocation != nil {
		pos := &Position{Latitude: float32(mvj.VehicleLocation.Latitude), Longitude: float32(mvj.VehicleLocation.Longitude)}
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
		if t, ok := ParseISOTime(*va.RecordedAtTime); ok {
			ts := t.Unix()
			vp.Timestamp = &ts
		}
	}
	ent.Vehicle = vp

	return &Entity{ID: id, Datasource: derefString(mvj.DataSource), Message: ent, TTL: ttl}
}

// ET -> TripUpdate
func MapETToTripUpdate(evj *EstimatedVehicleJourney, opts Options) *Entity {
	if evj == nil || evj.FramedVehicleJourneyRef == nil || evj.FramedVehicleJourneyRef.DatedVehicleJourneyRef == nil {
		return nil
	}

	tripId := *evj.FramedVehicleJourneyRef.DatedVehicleJourneyRef
	var startDate string
	if evj.OriginAimedDepartureTime != nil {
		if t, ok := ParseISOTime(*evj.OriginAimedDepartureTime); ok {
			startDate = FormatDateYYYYMMDD(t)
		}
	}

	id := tripId
	if startDate != "" {
		id = tripId + "-" + startDate
	}

	var latest time.Time
	for _, rc := range evj.RecordedCalls {
		times := []*string{rc.ActualArrivalTime, rc.ExpectedArrivalTime, rc.AimedArrivalTime, rc.ActualDepartureTime, rc.ExpectedDepartureTime, rc.AimedDepartureTime}
		for _, ts := range times {
			if ts != nil {
				if t, ok := ParseISOTime(*ts); ok {
					latest = Latest(latest, t)
				}
			}
		}
	}
	for _, ec := range evj.EstimatedCalls {
		times := []*string{ec.ExpectedArrivalTime, ec.AimedArrivalTime, ec.ExpectedDepartureTime, ec.AimedDepartureTime}
		for _, ts := range times {
			if ts != nil {
				if t, ok := ParseISOTime(*ts); ok {
					latest = Latest(latest, t)
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

	ent := &FeedEntity{Id: &id}
	tu := &TripUpdate{}
	td := &TripDescriptor{TripId: tripId}
	if evj.LineRef != nil {
		td.RouteId = *evj.LineRef
	}
	if evj.OriginAimedDepartureTime != nil {
		if t, ok := ParseISOTime(*evj.OriginAimedDepartureTime); ok {
			td.StartDate = FormatDateYYYYMMDD(t)
			td.StartTime = t.Format("15:04:05")
		}
	}
	tu.Trip = td
	if evj.VehicleRef != nil && *evj.VehicleRef != "" {
		tu.Vehicle = &VehicleDescriptor{Id: *evj.VehicleRef}
	}

	stopSeq := int32(0)
	for _, rc := range evj.RecordedCalls {
		stu := StopTimeUpdate{}
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
				stu.Arrival = &StopTimeEvent{Delay: d}
			}
		}
		if rc.AimedDepartureTime != nil {
			if d := delaySeconds(rc.AimedDepartureTime, rc.ActualDepartureTime, rc.ExpectedDepartureTime); d != nil {
				stu.Departure = &StopTimeEvent{Delay: d}
			}
		}
		tu.StopTimeUpdate = append(tu.StopTimeUpdate, stu)
		stopSeq++
	}
	for _, ec := range evj.EstimatedCalls {
		stu := StopTimeUpdate{}
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
				stu.Arrival = &StopTimeEvent{Delay: d}
			}
		}
		if ec.AimedDepartureTime != nil && ec.ExpectedDepartureTime != nil {
			if d := diffSeconds(*ec.AimedDepartureTime, *ec.ExpectedDepartureTime); d != nil {
				stu.Departure = &StopTimeEvent{Delay: d}
			}
		}
		tu.StopTimeUpdate = append(tu.StopTimeUpdate, stu)
		stopSeq++
	}

	ent.TripUpdate = tu

	return &Entity{
		ID:         id,
		Datasource: derefString(evj.DataSource),
		Message:    ent,
		TTL:        ttl,
	}
}

// SX -> Alert
func MapSXToAlert(sx *PtSituationElement, opts Options) *Entity {
	if sx == nil || sx.SituationNumber == nil {
		return nil
	}
	id := *sx.SituationNumber

	var end time.Time
	for _, vp := range sx.ValidityPeriods {
		if vp.EndTime != nil {
			if t, ok := ParseISOTime(*vp.EndTime); ok {
				if t.After(end) {
					end = t
				}
			}
		}
	}
	ttl := 365 * 24 * time.Hour
	if !end.IsZero() {
		if d := time.Until(end); d > 0 {
			ttl = d
		}
	}

	ent := &FeedEntity{Id: &id}
	alert := &Alert{}
	if len(sx.Summaries) > 0 {
		ts := TranslatedString{}
		for _, t := range sx.Summaries {
			if t.Value != "" {
				lang := t.Lang
				ts.Translation = append(ts.Translation, Translation{Text: t.Value, Language: strPtrOrNil(lang)})
			}
		}
		alert.HeaderText = &ts
	}
	if len(sx.Descriptions) > 0 {
		ts := TranslatedString{}
		for _, t := range sx.Descriptions {
			if t.Value != "" {
				lang := t.Lang
				ts.Translation = append(ts.Translation, Translation{Text: t.Value, Language: strPtrOrNil(lang)})
			}
		}
		alert.DescriptionText = &ts
	}
	for _, vp := range sx.ValidityPeriods {
		tr := TimeRange{}
		if vp.StartTime != nil {
			if t, ok := ParseISOTime(*vp.StartTime); ok {
				ts := t.Unix()
				tr.Start = &ts
			}
		}
		if vp.EndTime != nil {
			if t, ok := ParseISOTime(*vp.EndTime); ok {
				te := t.Unix()
				tr.End = &te
			}
		}
		alert.ActivePeriod = append(alert.ActivePeriod, tr)
	}
	if sx.Affects != nil {
		for _, sp := range sx.Affects.StopPoints {
			if sp.StopPointRef != nil {
				sid := *sp.StopPointRef
				alert.InformedEntity = append(alert.InformedEntity, EntitySelector{StopId: &sid})
			}
		}
		for _, vj := range sx.Affects.VehicleJourneys {
			if vj.LineRef != nil {
				rid := *vj.LineRef
				alert.InformedEntity = append(alert.InformedEntity, EntitySelector{RouteId: &rid})
			}
			if vj.FramedVehicleJourneyRef != nil && vj.FramedVehicleJourneyRef.DatedVehicleJourneyRef != nil {
				td := TripDescriptor{TripId: *vj.FramedVehicleJourneyRef.DatedVehicleJourneyRef}
				if vj.FramedVehicleJourneyRef.DataFrameRef != nil {
					td.StartDate = sanitizeDate(*vj.FramedVehicleJourneyRef.DataFrameRef)
				}
				alert.InformedEntity = append(alert.InformedEntity, EntitySelector{Trip: &td})
			}
			for _, dvj := range vj.DatedVehicleJourneyRefs {
				td := TripDescriptor{TripId: dvj}
				if vj.OriginAimedDepartureTime != nil {
					if t, ok := ParseISOTime(*vj.OriginAimedDepartureTime); ok {
						td.StartDate = FormatDateYYYYMMDD(t)
					}
				}
				alert.InformedEntity = append(alert.InformedEntity, EntitySelector{Trip: &td})
			}
			for _, r := range vj.Routes {
				for _, sp := range r.StopPoints.StopPoints {
					if sp.StopPointRef != nil {
						sid := *sp.StopPointRef
						alert.InformedEntity = append(alert.InformedEntity, EntitySelector{StopId: &sid})
					}
				}
			}
		}
		for _, net := range sx.Affects.Networks {
			for _, line := range net.AffectedLines {
				if line.LineRef != nil {
					rid := *line.LineRef
					alert.InformedEntity = append(alert.InformedEntity, EntitySelector{RouteId: &rid})
				}
				for _, r := range line.Routes {
					for _, sp := range r.StopPoints.StopPoints {
						if sp.StopPointRef != nil {
							sid := *sp.StopPointRef
							alert.InformedEntity = append(alert.InformedEntity, EntitySelector{RouteId: line.LineRef, StopId: &sid})
						}
					}
				}
			}
		}
	}
	if len(sx.InfoLinks) > 0 {
		ts := TranslatedString{}
		for _, l := range sx.InfoLinks {
			ts.Translation = append(ts.Translation, Translation{Text: l.Uri})
		}
		alert.Url = &ts
	}

	ent.Alert = alert
	return &Entity{ID: id, Datasource: derefString(sx.ParticipantRef), Message: ent, TTL: ttl}
}

// helpers
func derefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func sanitizeDate(s string) string {
	if len(s) == 10 && s[4] == '-' && s[7] == '-' {
		return s[0:4] + s[5:7] + s[8:10]
	}
	return s
}

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
	at, ok1 := ParseISOTime(aimed)
	ut, ok2 := ParseISOTime(updated)
	if !ok1 || !ok2 {
		return nil
	}
	d := int32(ut.Sub(at).Seconds())
	return &d
}
