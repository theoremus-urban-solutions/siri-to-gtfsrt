package convert

import (
	"golang/gtfsrt"
	"golang/mapper"
	"golang/siri"
	"golang/types"
)

func convertSIRI(sd *siri.ServiceDelivery, opts types.Options) ([]types.Entity, error) {
	var out []types.Entity
	if sd == nil {
		return out, nil
	}

	for _, d := range sd.EstimatedTimetableDeliveries {
		for _, f := range d.EstimatedJourneyVersionFrames {
			for _, evj := range f.EstimatedVehicleJourneys {
				if e := mapper.MapETToTripUpdate(&evj, opts); e != nil {
					e.Kind = "trip_update"
					out = append(out, *e)
				}
			}
		}
	}

	for _, d := range sd.VehicleMonitoringDeliveries {
		for _, va := range d.VehicleActivities {
			if e := mapper.MapVMToVehiclePosition(&va, opts); e != nil {
				e.Kind = "vehicle_position"
				out = append(out, *e)
			}
		}
	}

	for _, d := range sd.SituationExchangeDeliveries {
		for _, sx := range d.Situations {
			if e := mapper.MapSXToAlert(&sx, opts); e != nil {
				e.Kind = "alert"
				out = append(out, *e)
			}
		}
	}
	return out, nil
}

func buildFeedMessage(entities []types.Entity) *gtfsrt.FeedMessage {
	msg := gtfsrt.NewFeedMessage()
	for _, e := range entities {
		if e.Message != nil {
			msg.Entity = append(msg.Entity, e.Message)
		}
	}
	return msg
}

func buildPerDatasource(entities []types.Entity) map[string]*gtfsrt.FeedMessage {
	out := make(map[string]*gtfsrt.FeedMessage)
	for _, e := range entities {
		if e.Message == nil {
			continue
		}
		msg, ok := out[e.Datasource]
		if !ok {
			msg = gtfsrt.NewFeedMessage()
			out[e.Datasource] = msg
		}
		msg.Entity = append(msg.Entity, e.Message)
	}
	return out
}
