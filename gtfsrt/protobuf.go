package gtfsrt

import (
	"strings"
	"time"

	gtfs "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

// ToProto converts the internal FeedMessage to a GTFS-RT protobuf FeedMessage.
func ToProto(m *FeedMessage) *gtfs.FeedMessage {
	if m == nil {
		return nil
	}
	pm := &gtfs.FeedMessage{}
	if m.Header != nil {
		pm.Header = &gtfs.FeedHeader{}
		if m.Header.Timestamp != nil {
			pm.Header.Timestamp = m.Header.Timestamp
		} else {
			// Ensure a header timestamp exists for consumer compatibility
			ts := uint64(time.Now().Unix())
			pm.Header.Timestamp = &ts
		}
		if m.Header.GtfsRealtimeVersion != nil {
			pm.Header.GtfsRealtimeVersion = m.Header.GtfsRealtimeVersion
		}
		if m.Header.Incrementality != nil {
			inc := gtfs.FeedHeader_Incrementality(*m.Header.Incrementality)
			pm.Header.Incrementality = &inc
		}
	}
	for _, e := range m.Entity {
		pm.Entity = append(pm.Entity, toProtoEntity(e))
	}
	return pm
}

func toProtoEntity(e *FeedEntity) *gtfs.FeedEntity {
	pe := &gtfs.FeedEntity{}
	if e == nil {
		return pe
	}
	pe.Id = e.Id
	if e.TripUpdate != nil {
		pe.TripUpdate = toProtoTripUpdate(e.TripUpdate)
	}
	if e.Vehicle != nil {
		pe.Vehicle = toProtoVehicle(e.Vehicle)
	}
	if e.Alert != nil {
		pe.Alert = toProtoAlert(e.Alert)
	}
	return pe
}

func toProtoTripUpdate(tu *TripUpdate) *gtfs.TripUpdate {
	ptu := &gtfs.TripUpdate{}
	if tu.Trip != nil {
		ptu.Trip = &gtfs.TripDescriptor{
			TripId:    proto.String(tu.Trip.TripId),
			RouteId:   proto.String(tu.Trip.RouteId),
			StartDate: proto.String(tu.Trip.StartDate),
			StartTime: proto.String(tu.Trip.StartTime),
		}
	}
	if tu.Vehicle != nil {
		ptu.Vehicle = &gtfs.VehicleDescriptor{Id: proto.String(tu.Vehicle.Id)}
	}
	for _, stu := range tu.StopTimeUpdate {
		ps := &gtfs.TripUpdate_StopTimeUpdate{StopId: proto.String(stu.StopId), StopSequence: proto.Uint32(uint32(stu.StopSequence))}
		if stu.Arrival != nil && stu.Arrival.Delay != nil {
			ps.Arrival = &gtfs.TripUpdate_StopTimeEvent{Delay: proto.Int32(*stu.Arrival.Delay)}
		}
		if stu.Departure != nil && stu.Departure.Delay != nil {
			ps.Departure = &gtfs.TripUpdate_StopTimeEvent{Delay: proto.Int32(*stu.Departure.Delay)}
		}
		ptu.StopTimeUpdate = append(ptu.StopTimeUpdate, ps)
	}
	return ptu
}

func toProtoVehicle(v *VehiclePosition) *gtfs.VehiclePosition {
	pv := &gtfs.VehiclePosition{}
	if v.Trip != nil {
		pv.Trip = &gtfs.TripDescriptor{
			TripId:    proto.String(v.Trip.TripId),
			RouteId:   proto.String(v.Trip.RouteId),
			StartDate: proto.String(v.Trip.StartDate),
			StartTime: proto.String(v.Trip.StartTime),
		}
	}
	if v.Vehicle != nil {
		pv.Vehicle = &gtfs.VehicleDescriptor{Id: proto.String(v.Vehicle.Id)}
	}
	if v.Position != nil {
		pp := &gtfs.Position{Latitude: proto.Float32(v.Position.Latitude), Longitude: proto.Float32(v.Position.Longitude)}
		if v.Position.Bearing != nil {
			pp.Bearing = proto.Float32(*v.Position.Bearing)
		}
		if v.Position.Speed != nil {
			pp.Speed = proto.Float32(*v.Position.Speed)
		}
		if v.Position.Odometer != nil {
			pp.Odometer = proto.Float64(*v.Position.Odometer)
		}
		pv.Position = pp
	}
	if v.Timestamp != nil {
		ts := uint64(*v.Timestamp)
		pv.Timestamp = &ts
	}
	if st := mapVehicleStopStatus(v.CurrentStatus); st != nil {
		pv.CurrentStatus = st
	}
	if v.StopId != nil {
		pv.StopId = proto.String(*v.StopId)
	}
	if v.CurrentStopSequence != nil {
		pv.CurrentStopSequence = proto.Uint32(uint32(*v.CurrentStopSequence))
	}
	if v.OccupancyStatus != nil {
		if os := mapOccupancyStatus(*v.OccupancyStatus); os != nil {
			pv.OccupancyStatus = os
		}
	}
	if v.CongestionLevel != nil {
		if cl := mapCongestionLevel(*v.CongestionLevel); cl != nil {
			pv.CongestionLevel = cl
		}
	}
	return pv
}

func toProtoAlert(a *Alert) *gtfs.Alert {
	pa := &gtfs.Alert{}
	if a.HeaderText != nil {
		pa.HeaderText = toProtoTranslatedString(a.HeaderText)
	}
	if a.DescriptionText != nil {
		pa.DescriptionText = toProtoTranslatedString(a.DescriptionText)
	}
	// Map cause and effect; default when missing or unmapped
	{
		var causeStr string
		if a.Cause != nil {
			causeStr = *a.Cause
		}
		if c := mapAlertCause(causeStr); c != nil {
			pa.Cause = c
		} else {
			v := gtfs.Alert_UNKNOWN_CAUSE
			pa.Cause = &v
		}
	}
	{
		var effectStr string
		if a.Effect != nil {
			effectStr = *a.Effect
		}
		if e := mapAlertEffect(effectStr); e != nil {
			pa.Effect = e
		} else {
			v := gtfs.Alert_UNKNOWN_EFFECT
			pa.Effect = &v
		}
	}
	for _, tr := range a.ActivePeriod {
		ptr := &gtfs.TimeRange{}
		if tr.Start != nil {
			s := uint64(*tr.Start)
			ptr.Start = &s
		}
		if tr.End != nil {
			e := uint64(*tr.End)
			ptr.End = &e
		}
		pa.ActivePeriod = append(pa.ActivePeriod, ptr)
	}
	for _, ie := range a.InformedEntity {
		pie := &gtfs.EntitySelector{}
		if ie.RouteId != nil {
			pie.RouteId = proto.String(*ie.RouteId)
		}
		if ie.StopId != nil {
			pie.StopId = proto.String(*ie.StopId)
		}
		if ie.Trip != nil {
			pie.Trip = &gtfs.TripDescriptor{
				TripId:    proto.String(ie.Trip.TripId),
				RouteId:   proto.String(ie.Trip.RouteId),
				StartDate: proto.String(ie.Trip.StartDate),
				StartTime: proto.String(ie.Trip.StartTime),
			}
		}
		pa.InformedEntity = append(pa.InformedEntity, pie)
	}
	if a.Url != nil {
		pa.Url = toProtoTranslatedString(a.Url)
	}
	return pa
}

func toProtoTranslatedString(ts *TranslatedString) *gtfs.TranslatedString {
	pts := &gtfs.TranslatedString{}
	for _, tr := range ts.Translation {
		pt := &gtfs.TranslatedString_Translation{Text: proto.String(tr.Text)}
		if tr.Language != nil {
			pt.Language = proto.String(*tr.Language)
		}
		pts.Translation = append(pts.Translation, pt)
	}
	return pts
}

// MarshalPBF encodes the internal FeedMessage to PBF bytes.
func MarshalPBF(m *FeedMessage) ([]byte, error) {
	pm := ToProto(m)
	return proto.Marshal(pm)
}

// Mapping helpers from internal strings (from SIRI) to GTFS-RT enums
func mapVehicleStopStatus(s string) *gtfs.VehiclePosition_VehicleStopStatus {
	switch normalize(s) {
	case "incoming_at", "incoming", "approaching", "approach":
		v := gtfs.VehiclePosition_INCOMING_AT
		return &v
	case "stopped_at", "stopped", "at_stop":
		v := gtfs.VehiclePosition_STOPPED_AT
		return &v
	case "in_transit_to", "in_transit", "moving", "enroute":
		v := gtfs.VehiclePosition_IN_TRANSIT_TO
		return &v
	default:
		return nil
	}
}

func mapOccupancyStatus(s string) *gtfs.VehiclePosition_OccupancyStatus {
	switch normalize(s) {
	case "empty":
		v := gtfs.VehiclePosition_EMPTY
		return &v
	case "many_seats_available", "many":
		v := gtfs.VehiclePosition_MANY_SEATS_AVAILABLE
		return &v
	case "few_seats_available", "few":
		v := gtfs.VehiclePosition_FEW_SEATS_AVAILABLE
		return &v
	case "standing_room_only", "standing":
		v := gtfs.VehiclePosition_STANDING_ROOM_ONLY
		return &v
	case "crushed_standing_room_only", "crushed":
		v := gtfs.VehiclePosition_CRUSHED_STANDING_ROOM_ONLY
		return &v
	case "full":
		v := gtfs.VehiclePosition_FULL
		return &v
	case "not_accepting_passengers", "closed":
		v := gtfs.VehiclePosition_NOT_ACCEPTING_PASSENGERS
		return &v
	default:
		return nil
	}
}

func mapCongestionLevel(s string) *gtfs.VehiclePosition_CongestionLevel {
	switch normalize(s) {
	case "running_smoothly", "smooth":
		v := gtfs.VehiclePosition_RUNNING_SMOOTHLY
		return &v
	case "stop_and_go":
		v := gtfs.VehiclePosition_STOP_AND_GO
		return &v
	case "congestion":
		v := gtfs.VehiclePosition_CONGESTION
		return &v
	case "severe_congestion", "severe":
		v := gtfs.VehiclePosition_SEVERE_CONGESTION
		return &v
	case "unknown", "":
		v := gtfs.VehiclePosition_UNKNOWN_CONGESTION_LEVEL
		return &v
	default:
		return nil
	}
}

func normalize(s string) string {
	// Lowercase and trim spaces; SIRI values are usually simple tokens
	return strings.TrimSpace(strings.ToLower(s))
}

// Alert Cause mapping (case-insensitive) with transitives
func mapAlertCause(s string) *gtfs.Alert_Cause {
	switch normalize(s) {
	case "unknown", "unknown_cause":
		v := gtfs.Alert_UNKNOWN_CAUSE
		return &v
	case "other", "other_cause":
		v := gtfs.Alert_OTHER_CAUSE
		return &v
	case "technical_problem", "technical", "signalfailure", "signal_failure":
		v := gtfs.Alert_TECHNICAL_PROBLEM
		return &v
	case "strike", "staffunavailable", "staff_unavailable":
		v := gtfs.Alert_STRIKE
		return &v
	case "demonstration":
		v := gtfs.Alert_DEMONSTRATION
		return &v
	case "accident":
		v := gtfs.Alert_ACCIDENT
		return &v
	case "holiday":
		v := gtfs.Alert_HOLIDAY
		return &v
	case "weather":
		v := gtfs.Alert_WEATHER
		return &v
	case "maintenance", "maintenance_work":
		v := gtfs.Alert_MAINTENANCE
		return &v
	case "construction", "constructionwork", "construction_work":
		v := gtfs.Alert_CONSTRUCTION
		return &v
	case "police_activity", "securityalert", "security_alert":
		v := gtfs.Alert_POLICE_ACTIVITY
		return &v
	case "medical_emergency":
		v := gtfs.Alert_MEDICAL_EMERGENCY
		return &v
	default:
		// fallback to UNKNOWN_CAUSE as UNDEFINED is not in protobuf enum
		v := gtfs.Alert_UNKNOWN_CAUSE
		return &v
	}
}

// Alert Effect mapping (case-insensitive)
func mapAlertEffect(s string) *gtfs.Alert_Effect {
	switch normalize(s) {
	case "no_service":
		v := gtfs.Alert_NO_SERVICE
		return &v
	case "reduced_service":
		v := gtfs.Alert_REDUCED_SERVICE
		return &v
	case "significant_delays", "significantdelays":
		v := gtfs.Alert_SIGNIFICANT_DELAYS
		return &v
	case "detour":
		v := gtfs.Alert_DETOUR
		return &v
	case "additional_service":
		v := gtfs.Alert_ADDITIONAL_SERVICE
		return &v
	case "modified_service":
		v := gtfs.Alert_MODIFIED_SERVICE
		return &v
	case "stop_moved", "stoppointrelocation", "stop_point_relocation":
		v := gtfs.Alert_STOP_MOVED
		return &v
	case "other", "other_effect":
		v := gtfs.Alert_OTHER_EFFECT
		return &v
	case "unknown", "unknown_effect":
		v := gtfs.Alert_UNKNOWN_EFFECT
		return &v
	default:
		v := gtfs.Alert_UNKNOWN_EFFECT
		return &v
	}
}
