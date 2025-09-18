package siritogtfs

// Convert and feed building logic

func ConvertSIRI(sd *ServiceDelivery, opts Options) ([]Entity, error) {
	return convertSIRI(sd, opts)
}

func BuildFeedMessage(entities []Entity) *FeedMessage {
	return buildFeedMessage(entities)
}

func BuildPerDatasource(entities []Entity) map[string]*FeedMessage {
	return buildPerDatasource(entities)
}

func convertSIRI(sd *ServiceDelivery, opts Options) ([]Entity, error) {
	var out []Entity
	if sd == nil {
		return out, nil
	}

	for _, d := range sd.EstimatedTimetableDeliveries {
		for _, f := range d.EstimatedJourneyVersionFrames {
			for _, evj := range f.EstimatedVehicleJourneys {
				if e := MapETToTripUpdate(&evj, opts); e != nil {
					e.Kind = "trip_update"
					out = append(out, *e)
				}
			}
		}
	}

	for _, d := range sd.VehicleMonitoringDeliveries {
		for _, va := range d.VehicleActivities {
			if e := MapVMToVehiclePosition(&va, opts); e != nil {
				e.Kind = "vehicle_position"
				out = append(out, *e)
			}
		}
	}

	for _, d := range sd.SituationExchangeDeliveries {
		for _, sx := range d.Situations {
			if e := MapSXToAlert(&sx, opts); e != nil {
				e.Kind = "alert"
				out = append(out, *e)
			}
		}
	}
	return out, nil
}

func buildFeedMessage(entities []Entity) *FeedMessage {
	msg := NewFeedMessage()
	for _, e := range entities {
		if e.Message != nil {
			msg.Entity = append(msg.Entity, e.Message)
		}
	}
	return msg
}

func buildPerDatasource(entities []Entity) map[string]*FeedMessage {
	out := make(map[string]*FeedMessage)
	for _, e := range entities {
		if e.Message == nil {
			continue
		}
		msg, ok := out[e.Datasource]
		if !ok {
			msg = NewFeedMessage()
			out[e.Datasource] = msg
		}
		msg.Entity = append(msg.Entity, e.Message)
	}
	return out
}
