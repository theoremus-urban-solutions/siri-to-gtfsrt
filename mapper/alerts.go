package mapper

import (
	"golang/gtfsrt"
	"golang/siri"
	"golang/types"
	"time"
)

// MapSXToAlert converts a SIRI PtSituationElement into a GTFS-RT entity.
func MapSXToAlert(sx *siri.PtSituationElement, opts types.Options) *types.Entity {
	if sx == nil || sx.SituationNumber == nil {
		return nil
	}
	id := *sx.SituationNumber

	// TTL from validity periods, else 365d
	var end time.Time
	for _, vp := range sx.ValidityPeriods {
		if vp.EndTime != nil {
			if t, ok := siri.ParseISOTime(*vp.EndTime); ok {
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

	// Build Alert
	ent := &gtfsrt.FeedEntity{Id: &id}
	alert := &gtfsrt.Alert{}
	// Header/Description
	if len(sx.Summaries) > 0 {
		ts := gtfsrt.TranslatedString{}
		for _, t := range sx.Summaries {
			if t.Value != "" {
				lang := t.Lang
				ts.Translation = append(ts.Translation, gtfsrt.Translation{Text: t.Value, Language: strPtrOrNil(lang)})
			}
		}
		alert.HeaderText = &ts
	}
	if len(sx.Descriptions) > 0 {
		ts := gtfsrt.TranslatedString{}
		for _, t := range sx.Descriptions {
			if t.Value != "" {
				lang := t.Lang
				ts.Translation = append(ts.Translation, gtfsrt.Translation{Text: t.Value, Language: strPtrOrNil(lang)})
			}
		}
		alert.DescriptionText = &ts
	}
	// Active periods
	for _, vp := range sx.ValidityPeriods {
		tr := gtfsrt.TimeRange{}
		if vp.StartTime != nil {
			if t, ok := siri.ParseISOTime(*vp.StartTime); ok {
				ts := t.Unix()
				tr.Start = &ts
			}
		}
		if vp.EndTime != nil {
			if t, ok := siri.ParseISOTime(*vp.EndTime); ok {
				te := t.Unix()
				tr.End = &te
			}
		}
		alert.ActivePeriod = append(alert.ActivePeriod, tr)
	}
	// Informed entities
	if sx.Affects != nil {
		for _, sp := range sx.Affects.StopPoints {
			if sp.StopPointRef != nil {
				sid := *sp.StopPointRef
				alert.InformedEntity = append(alert.InformedEntity, gtfsrt.EntitySelector{StopId: &sid})
			}
		}
		for _, vj := range sx.Affects.VehicleJourneys {
			if vj.LineRef != nil {
				rid := *vj.LineRef
				alert.InformedEntity = append(alert.InformedEntity, gtfsrt.EntitySelector{RouteId: &rid})
			}
			if vj.FramedVehicleJourneyRef != nil && vj.FramedVehicleJourneyRef.DatedVehicleJourneyRef != nil {
				td := gtfsrt.TripDescriptor{TripId: *vj.FramedVehicleJourneyRef.DatedVehicleJourneyRef}
				if vj.FramedVehicleJourneyRef.DataFrameRef != nil {
					td.StartDate = sanitizeDate(*vj.FramedVehicleJourneyRef.DataFrameRef)
				}
				alert.InformedEntity = append(alert.InformedEntity, gtfsrt.EntitySelector{Trip: &td})
			}
			for _, dvj := range vj.DatedVehicleJourneyRefs {
				td := gtfsrt.TripDescriptor{TripId: dvj}
				if vj.OriginAimedDepartureTime != nil {
					if t, ok := siri.ParseISOTime(*vj.OriginAimedDepartureTime); ok {
						td.StartDate = siri.FormatDateYYYYMMDD(t)
					}
				}
				alert.InformedEntity = append(alert.InformedEntity, gtfsrt.EntitySelector{Trip: &td})
			}
			for _, r := range vj.Routes {
				for _, sp := range r.StopPoints.StopPoints {
					if sp.StopPointRef != nil {
						sid := *sp.StopPointRef
						alert.InformedEntity = append(alert.InformedEntity, gtfsrt.EntitySelector{StopId: &sid})
					}
				}
			}
		}
		for _, net := range sx.Affects.Networks {
			for _, line := range net.AffectedLines {
				if line.LineRef != nil {
					rid := *line.LineRef
					alert.InformedEntity = append(alert.InformedEntity, gtfsrt.EntitySelector{RouteId: &rid})
				}
				for _, r := range line.Routes {
					for _, sp := range r.StopPoints.StopPoints {
						if sp.StopPointRef != nil {
							sid := *sp.StopPointRef
							alert.InformedEntity = append(alert.InformedEntity, gtfsrt.EntitySelector{RouteId: line.LineRef, StopId: &sid})
						}
					}
				}
			}
		}
	}
	// URLs
	if len(sx.InfoLinks) > 0 {
		ts := gtfsrt.TranslatedString{}
		for _, l := range sx.InfoLinks {
			ts.Translation = append(ts.Translation, gtfsrt.Translation{Text: l.Uri})
		}
		alert.Url = &ts
	}

	ent.Alert = alert
	return &types.Entity{ID: id, Datasource: derefStrPtr(sx.ParticipantRef), Message: ent, TTL: ttl}
}

func derefStrPtr(p *string) string {
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
