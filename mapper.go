package siritogtfs

import (
	"strings"
	"time"
)

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
	// Strip SOFIA:SituationNumber: prefix
	id := stripPrefix(*sx.SituationNumber, "SOFIA:SituationNumber:")

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

	// Parse cause and effect from Summary text
	var bgSummary, enSummary string
	if len(sx.Summaries) > 0 {
		for _, t := range sx.Summaries {
			lang := t.Lang
			if lang == "" {
				lang = "en" // Default to English if no language specified
			}
			if lang == "bg" {
				bgSummary = t.Value
			} else if lang == "en" {
				enSummary = t.Value
			}
		}
	}

	// Use English summary to parse cause and effect, fallback to any summary
	summaryForParsing := enSummary
	if summaryForParsing == "" && len(sx.Summaries) > 0 {
		summaryForParsing = sx.Summaries[0].Value
	}

	if summaryForParsing != "" {
		cause := mapCauseIntToString(parseCauseFromSummary(summaryForParsing))
		alert.Cause = &cause
		effect := mapEffectIntToString(parseEffectFromSummary(summaryForParsing))
		alert.Effect = &effect
	} else {
		// Default values if no summary
		defaultCause := "OTHER_CAUSE"
		defaultEffect := "OTHER_EFFECT"
		alert.Cause = &defaultCause
		alert.Effect = &defaultEffect
	}

	if sx.Severity != nil {
		s := *sx.Severity
		alert.Severity = &s
	}
	alert.HeaderText = &TranslatedString{
		Translation: []Translation{
			{Text: bgSummary, Language: strPtr("bg")},
			{Text: enSummary, Language: strPtr("en")},
		},
	}
	var bgDescription, enDescription string
	if len(sx.Descriptions) > 0 {
		for _, t := range sx.Descriptions {
			lang := t.Lang
			if lang == "" {
				lang = "en" // Default to English if no language specified
			}
			if lang == "bg" {
				bgDescription = t.Value
			} else if lang == "en" {
				enDescription = t.Value
			}
		}
	}
	alert.DescriptionText = &TranslatedString{
		Translation: []Translation{
			{Text: bgDescription, Language: strPtr("bg")},
			{Text: enDescription, Language: strPtr("en")},
		},
	}
	// Prefer PublicationWindow if present; else fallback to ValidityPeriods
	if sx.PublicationWindow != nil && (sx.PublicationWindow.StartTime != nil || sx.PublicationWindow.EndTime != nil) {
		tr := TimeRange{}
		if sx.PublicationWindow.StartTime != nil {
			if t, ok := ParseISOTime(*sx.PublicationWindow.StartTime); ok {
				ts := t.Unix()
				tr.Start = &ts
			}
		}
		if sx.PublicationWindow.EndTime != nil {
			if t, ok := ParseISOTime(*sx.PublicationWindow.EndTime); ok {
				te := t.Unix()
				tr.End = &te
			}
		}
		alert.ActivePeriod = append(alert.ActivePeriod, tr)
	} else {
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
	}
	if sx.Affects != nil {
		for _, sp := range sx.Affects.StopPoints {
			if sp.StopPointRef != nil {
				sid := stripPrefix(*sp.StopPointRef, "SOFIA:Quay:")
				alert.InformedEntity = append(alert.InformedEntity, EntitySelector{StopId: &sid})
			}
		}
		for _, vj := range sx.Affects.VehicleJourneys {
			if vj.LineRef != nil {
				rid := stripPrefix(*vj.LineRef, "SOFIA:Line:")
				alert.InformedEntity = append(alert.InformedEntity, EntitySelector{RouteId: &rid})
			}
			if vj.FramedVehicleJourneyRef != nil && vj.FramedVehicleJourneyRef.DatedVehicleJourneyRef != nil {
				tid := stripPrefix(*vj.FramedVehicleJourneyRef.DatedVehicleJourneyRef, "SOFIA:ServiceJourney:")
				td := TripDescriptor{TripId: tid}
				if vj.FramedVehicleJourneyRef.DataFrameRef != nil {
					td.StartDate = sanitizeDate(*vj.FramedVehicleJourneyRef.DataFrameRef)
				}
				alert.InformedEntity = append(alert.InformedEntity, EntitySelector{Trip: &td})
			}
			for _, dvj := range vj.DatedVehicleJourneyRefs {
				tid := stripPrefix(dvj, "SOFIA:ServiceJourney:")
				td := TripDescriptor{TripId: tid}
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
						sid := stripPrefix(*sp.StopPointRef, "SOFIA:Quay:")
						alert.InformedEntity = append(alert.InformedEntity, EntitySelector{StopId: &sid})
					}
				}
			}
		}
		for _, net := range sx.Affects.Networks {
			for _, line := range net.AffectedLines {
				if line.LineRef != nil {
					rid := stripPrefix(*line.LineRef, "SOFIA:Line:")
					alert.InformedEntity = append(alert.InformedEntity, EntitySelector{RouteId: &rid})
				}
				for _, r := range line.Routes {
					for _, sp := range r.StopPoints.StopPoints {
						if sp.StopPointRef != nil {
							sid := stripPrefix(*sp.StopPointRef, "SOFIA:Quay:")
							ridPtr := line.LineRef
							if ridPtr != nil {
								rid := stripPrefix(*ridPtr, "SOFIA:Line:")
								alert.InformedEntity = append(alert.InformedEntity, EntitySelector{RouteId: &rid, StopId: &sid})
							} else {
								alert.InformedEntity = append(alert.InformedEntity, EntitySelector{StopId: &sid})
							}
						}
					}
				}
			}
		}
	}
	var bgUrl, enUrl string
	if len(sx.InfoLinks) > 0 {
		for _, l := range sx.InfoLinks {
			lang := l.Lang
			if lang == "" {
				lang = "en" // Default to English if no language specified
			}
			if lang == "bg" {
				bgUrl = l.Uri
			} else if lang == "en" {
				enUrl = l.Uri
			}
		}
	}
	alert.Url = &TranslatedString{
		Translation: []Translation{
			{Text: bgUrl, Language: strPtr("bg")},
			{Text: enUrl, Language: strPtr("en")},
		},
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

func strPtr(s string) *string {
	return &s
}

func stripPrefix(s, prefix string) string {
	if len(s) > len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}

func mapCauseIntToString(cause int32) string {
	switch cause {
	case 1:
		return "UNKNOWN_CAUSE"
	case 2:
		return "OTHER_CAUSE"
	case 3:
		return "TECHNICAL_PROBLEM"
	case 4:
		return "STRIKE"
	case 5:
		return "DEMONSTRATION"
	case 6:
		return "ACCIDENT"
	case 7:
		return "HOLIDAY"
	case 8:
		return "WEATHER"
	case 9:
		return "MAINTENANCE"
	case 10:
		return "CONSTRUCTION"
	case 11:
		return "POLICE_ACTIVITY"
	case 12:
		return "MEDICAL_EMERGENCY"
	default:
		return "UNKNOWN_CAUSE"
	}
}

func mapEffectIntToString(effect int32) string {
	switch effect {
	case 1:
		return "NO_SERVICE"
	case 2:
		return "REDUCED_SERVICE"
	case 3:
		return "SIGNIFICANT_DELAYS"
	case 4:
		return "DETOUR"
	case 5:
		return "ADDITIONAL_SERVICE"
	case 6:
		return "MODIFIED_SERVICE"
	case 7:
		return "OTHER_EFFECT"
	case 8:
		return "UNKNOWN_EFFECT"
	case 9:
		return "STOP_MOVED"
	case 10:
		return "NO_EFFECT"
	case 11:
		return "ACCESSIBILITY_ISSUE"
	default:
		return "UNKNOWN_EFFECT"
	}
}

func parseCauseFromSummary(summary string) int32 {
	// Extract cause part from "X:Y" format (e.g., "Maintenance:Stop moved")
	s := strings.ToLower(summary)

	// Try to extract just the cause portion (before the colon)
	causePart := s
	if idx := strings.Index(s, ":"); idx > 0 {
		causePart = s[:idx]
	}

	if strings.Contains(causePart, "maintenance") || strings.Contains(causePart, "поддръжка") {
		return 9 // MAINTENANCE
	}
	if strings.Contains(causePart, "construction") || strings.Contains(causePart, "строителна") {
		return 10 // CONSTRUCTION
	}
	if strings.Contains(causePart, "technical problem") || strings.Contains(causePart, "технически проблем") {
		return 3 // TECHNICAL_PROBLEM
	}
	if strings.Contains(causePart, "strike") || strings.Contains(causePart, "стачка") {
		return 4 // STRIKE
	}
	if strings.Contains(causePart, "demonstration") || strings.Contains(causePart, "демонстрация") {
		return 5 // DEMONSTRATION
	}
	if strings.Contains(causePart, "accident") || strings.Contains(causePart, "авария") {
		return 6 // ACCIDENT
	}
	if strings.Contains(causePart, "holiday") || strings.Contains(causePart, "праздник") {
		return 7 // HOLIDAY
	}
	if strings.Contains(causePart, "weather") || strings.Contains(causePart, "време") {
		return 8 // WEATHER
	}
	if strings.Contains(causePart, "police") || strings.Contains(causePart, "полиц") {
		return 11 // POLICE_ACTIVITY
	}
	if strings.Contains(causePart, "medical") || strings.Contains(causePart, "медицин") {
		return 12 // MEDICAL_EMERGENCY
	}
	if strings.Contains(causePart, "unknown") || strings.Contains(causePart, "неизвестно") {
		return 1 // UNKNOWN_CAUSE
	}
	if strings.Contains(causePart, "other") || strings.Contains(causePart, "друго") {
		return 2 // OTHER_CAUSE
	}

	return 2 // OTHER_CAUSE (default)
}

func parseEffectFromSummary(summary string) int32 {
	// Extract effect part from "X:Y" format (e.g., "Maintenance:Stop moved")
	s := strings.ToLower(summary)

	// Try to extract just the effect portion (after the colon)
	effectPart := s
	if idx := strings.Index(s, ":"); idx > 0 && idx < len(s)-1 {
		effectPart = s[idx+1:]
	}

	if strings.Contains(effectPart, "no service") || strings.Contains(effectPart, "не се изпълнява") {
		return 1 // NO_SERVICE
	}
	if strings.Contains(effectPart, "reduced service") || strings.Contains(effectPart, "понижено обслужване") {
		return 2 // REDUCED_SERVICE
	}
	if strings.Contains(effectPart, "significant delay") || strings.Contains(effectPart, "значителни закъснения") {
		return 3 // SIGNIFICANT_DELAYS
	}
	if strings.Contains(effectPart, "detour") || strings.Contains(effectPart, "отклонение") {
		return 4 // DETOUR
	}
	if strings.Contains(effectPart, "additional service") || strings.Contains(effectPart, "допълнително обслужване") {
		return 5 // ADDITIONAL_SERVICE
	}
	if strings.Contains(effectPart, "modified service") || strings.Contains(effectPart, "модифицирано обслужване") {
		return 6 // MODIFIED_SERVICE
	}
	if strings.Contains(effectPart, "stop moved") || strings.Contains(effectPart, "преместена спирка") {
		return 9 // STOP_MOVED
	}
	if strings.Contains(effectPart, "no impact") || strings.Contains(effectPart, "no effect") || strings.Contains(effectPart, "няма ефект") {
		return 10 // NO_EFFECT
	}
	if strings.Contains(effectPart, "accessibility") || strings.Contains(effectPart, "достъпност") {
		return 11 // ACCESSIBILITY_ISSUE
	}
	if strings.Contains(effectPart, "unknown") || strings.Contains(effectPart, "неизвестно") {
		return 8 // UNKNOWN_EFFECT
	}
	if strings.Contains(effectPart, "other") || strings.Contains(effectPart, "друго") {
		return 7 // OTHER_EFFECT
	}

	return 7 // OTHER_EFFECT (default)
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
