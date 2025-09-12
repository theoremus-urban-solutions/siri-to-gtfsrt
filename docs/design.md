### Go Library Design — SIRI (ET/VM/SX) → GTFS-Realtime

**Scope**: Provide a Go library and helpers to convert SIRI realtime data into GTFS-RT. No servers, no Pub/Sub; pure code + a thin CLI. This mirrors the Java logic in `SiriToGtfsRealtimeService`, `GtfsRtMapper`, and `AlertFactory`.

### Inputs and outputs
- **Input**: SIRI data as JSON decoded into Go structs (compatible with Entur's Avro model fields). Initial focus is the already-merged ServiceDelivery types: ET, VM, SX. XML support can be added later.
- **Output**: GTFS-RT `FeedMessage` (proto) and/or individual `FeedEntity` objects. JSON output supported via proto JSON marshaling when requested.

### Package layout (library only)
- `pkg/siri`:
  - Go structs for the subset of SIRI fields we need (ET/VM/SX). Designed to unmarshal JSON matching the existing Avro-shaped payloads.
- `pkg/gtfsrt`:
  - Protobuf-generated Go code (`google.golang.org/protobuf`), and small helpers to create headers/messages.
- `pkg/convert`:
  - Core conversion logic ET→TripUpdates, VM→VehiclePositions, SX→Alerts.
  - TTL computation and ID/key derivation.
- `pkg/types`:
  - Common types used across converters (e.g., `Entity`, `Options`).
- `pkg/enrich` (optional, future):
  - Interface for service-journey enrichment (GraphQL client plug-in).

### Public API (draft)
```go
package convert

import (
    "time"
    gtfsrt "github.com/google/transit/gtfs-realtime/go" // canonical module path TBD
)

type Options struct {
    // Filters
    ETWhitelist []string
    VMWhitelist []string
    SXWhitelist []string

    // Vehicle proximity heuristics (VM)
    CloseToNextStopPercentage int // default 95
    CloseToNextStopDistance   int // meters, default 500

    // TTL fallbacks
    VMGracePeriod time.Duration // default 5m
}

type Entity struct {
    ID         string
    Datasource string
    Message    *gtfsrt.FeedEntity
    TTL        time.Duration
}

// Convert a whole SIRI ServiceDelivery payload (if present) aggregating ET/VM/SX.
func ConvertSIRI(sd *siri.ServiceDelivery, opts Options) ([]Entity, error)

// Focused converters for each SIRI branch
func ConvertET(et *siri.EstimatedVehicleJourney, opts Options) ([]Entity, error)
func ConvertVM(vm *siri.VehicleActivity, opts Options) ([]Entity, error)
func ConvertSX(sx *siri.PtSituationElement, opts Options) ([]Entity, error)

// Message builders
func BuildFeedMessage(entities []Entity) *gtfsrt.FeedMessage
func BuildPerDatasource(entities []Entity) map[string]*gtfsrt.FeedMessage
```

Notes:
- `ConvertSIRI` walks `EstimatedTimetableDeliveries`, `VehicleMonitoringDeliveries`, `SituationExchangeDeliveries` and merges results.
- Whitelists: if defined and a record’s datasource/participant is not in the list, skip it (and return no error).
- `Entity.ID`:
  - ET: `TripDescriptor(trip_id,start_date) [+ vehicle_id if present]` (mirrors Java keying)
  - VM: `vehicle_id` if present else `trip_id-start_date`
  - SX: `situation_number`
- `Entity.TTL`:
  - ET: max of actual/expected arrival/departure times across recorded/estimated calls; fallback to `VMGracePeriod`
  - VM: `validUntilTime` if present else `VMGracePeriod`
  - SX: end of validity periods; fallback 365d

### Conversion rules (parity with Java)
- TripUpdates (ET):
  - TripDescriptor: from `FramedVehicleJourneyRef.datedVehicleJourneyRef` or resolved `DatedVehicleJourneyRef` via enrichment (optional).
  - RouteId: from `lineRef` if present.
  - StartDate/StartTime: from `originAimedDepartureTime` or enrichment date.
  - StopTimeUpdates: compute arrival/departure delays from aimed vs actual/expected; stopSequence uses `order-1` (or incremental fallback).
- VehiclePositions (VM):
  - TripDescriptor: from `FramedVehicleJourneyRef` or resolved `VehicleJourneyRef`.
  - VehicleDescriptor: from `vehicleRef` when present.
  - Position: lat/lon, bearing, speed (velocity), timestamp (recordedAtTime).
  - CurrentStatus: `STOPPED_AT` if `vehicleAtStop`, else `INCOMING_AT` if close-to-next-stop, else `IN_TRANSIT_TO`.
  - Odometer: derived from `ProgressBetweenStops` (`percentage * linkDistance`).
  - Occupancy, congestion mapped to GTFS-RT enums where possible.
- Alerts (SX):
  - Header/Description: from translated strings.
  - ActivePeriod: from validity periods with fallback end.
  - InformedEntity: from stop points, routes/lines, trips (via multiple references), and combination of trip+stop when available.

### Error handling
- Best-effort conversion: invalid/missing fields skip entity creation without failing the whole batch.
- Return errors only for critical failures (e.g., JSON decode issues) in CLI adapters; core converters favor empty results and nil errors.

### Extensibility
- Add XML decoder in `pkg/siri/xml` implementing the same structs.
- Plug-in enrichment via an interface (GraphQL or other sources) to resolve `DatedServiceJourney`.

### Non-goals (for now)
- No HTTP servers, no Pub/Sub clients, no Redis. Those belong to the future hub service.


