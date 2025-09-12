### TODO — Port SIRI → GTFS-RT Library (reuse-first plan)

Goal: Build a Go library (plus thin CLI) that converts SIRI ET/VM/SX into GTFS-RT with maximum reuse of existing logic from this repo.

### Reusable logic to port (priority order)
- Siri→GTFS-RT conversion orchestration:
  - `SiriToGtfsRealtimeService.convertSiriToGtfsRt` (walk ServiceDelivery; merge ET/VM/SX)
  - `convertSiriEtToGtfsRt`, `convertSiriVmToGtfsRt`, `convertSiriSxToGtfsRt` (focused converters)
- Builders and helpers:
  - `GtfsRealtimeLibrary.createFeedMessageBuilder` (header init and timestamp)
  - `GtfsRtMapper` (core field-level mapping for TripUpdates and VehiclePositions)
  - `AlertFactory` (Situation → GTFS-RT Alert mapping including informed entities)
  - `AvroHelper.translation`, `AvroHelper.getInstant` (translated strings, timestamp parsing)
  - `SiriLibrary.getLatestTimestamp` (latest-of-two instants)
- ID and TTL rules:
  - Trip/vehicle keying from `getTripIdForEstimatedVehicleJourney`, `getVehicleIdForKey`
  - TTL from `validUntilTime`, call times, validity periods, with grace/fallbacks
- Whitelisting and metrics behavior (port filters; metrics hooks as no-op initially)
- Tests as specifications (convert to Go test fixtures):
  - `TestSiriVMToGtfsRealtimeService` (status/occupancy/odometer, filtering)
  - `TestSiriETToGtfsRealtimeService` (stop-time delays, sequencing, filtering)
  - `TestSiriSXToGtfsRealtimeService` and `TestAlertFactory` (alert content)

### Items to drop (for this library)
- Spring Boot/Camel routes, HTTP server, Pub/Sub, Redis storage, Prometheus registry (leave stubs/hooks in design for future use).
- GraphQL enrichment (make a pluggable interface; default no-op).

### Data model equivalence (Go)
- Define Go structs mirroring the used subset of SIRI Avro model fields (ET/VM/SX). JSON-only initially.
- Use official GTFS-RT Go protobufs for output.

### Step-by-step plan (implementation)
1) Module scaffold
   - Create `go.mod` and directories: `pkg/siri`, `pkg/gtfsrt`, `pkg/convert`, `cmd/kishar-convert`.
   - Add dependencies: `google.golang.org/protobuf`, `github.com/golang/protobuf/ptypes` (or modern), `github.com/stretchr/testify` for tests.

2) GTFS-RT helpers
   - Implement `pkg/gtfsrt/header.go`: `NewFeedMessageBuilder()` to mirror Java header (timestamp, FULL_DATASET, version "1.0").

3) SIRI structs (subset)
   - Implement minimal JSON structs for:
     - ServiceDelivery with `EstimatedTimetableDeliveries`, `VehicleMonitoringDeliveries`, `SituationExchangeDeliveries`.
     - ET: `EstimatedVehicleJourney` with framed refs, calls, recorded calls, lineRef, vehicleRef, dataSource, originAimedDepartureTime.
     - VM: `VehicleActivity`: `RecordedAtTime`, `ValidUntilTime`, `MonitoredVehicleJourney` (vehicleRef/lineRef/vehicleLocation/bearing/velocity/dataSource/monitoredCall), `ProgressBetweenStops`.
     - SX: `PtSituationElement`: `ParticipantRef`, `SituationNumber`, `ValidityPeriods`, `Summaries`, `Descriptions`, `Affects` (stops/routes/trips), `InfoLinks`.
   - Provide `time.Time` parsing for ISO-8601; helpers for "latest timestamp".

4) Converters (pkg/convert)
   - Options and Entity types (whitelists, proximity thresholds, grace periods).
   - ET → TripUpdate:
     - TripDescriptor from framed/dvj ref; routeId from lineRef; start date/time from originAimedDepartureTime.
     - StopTimeUpdates with delays from aimed vs actual/expected; stopSequence from order-1.
     - TTL = max of call times; fallback to grace.
   - VM → VehiclePosition:
     - TripDescriptor from framed ref or resolved dvj.
     - VehicleDescriptor from vehicleRef; position lat/lon/bearing/speed; timestamp recordedAtTime.
     - CurrentStatus: STOPPED_AT | INCOMING_AT | IN_TRANSIT_TO (proximity rules).
     - Occupancy/congestion mapping.
     - TTL = validUntilTime or grace; entity ID selection per Java.
   - SX → Alert:
     - Header/description, active periods with fallback end; informed entities from stops, routes, trips; URLs.

5) Aggregation helpers
   - `BuildFeedMessage(entities []Entity)` and `BuildPerDatasource(...)` to create messages and maps.

6) CLI (`cmd/kishar-convert`)
   - Command `siri2gtfs` with flags from `cli.md`.
   - Read JSON from stdin/file/url, decode, call `convert.ConvertSIRI`, output GTFS-RT as PBF or JSON.

7) Tests (mirror Java tests)
   - Create Go fixtures using the XML snippets in tests via a tiny XML→JSON conversion step or direct JSON fixtures equivalent to Avro-shaped JSON.
   - VM tests: status transitions, occupancy mapping, proximity rules, datasource filtering.
   - ET tests: delays, stop sequence, route id, filtering.
   - SX tests: header/description, active period, empty informed entities path, filtering.

8) Documentation
   - Update `design.md` as the API stabilizes; document struct fields required from SIRI JSON.
   - CLI usage examples verified with real outputs.

9) Future hooks
   - Define optional interfaces for enrichment and metrics (no-op default) to keep drop-in parity for the hub.

### Notes for faithful parity
- Match Java date formats for GTFS-RT `start_date` (YYYYMMDD) and `start_time` (HH:MM:SS).
- Preserve entity ID construction and TTL rules to ensure deterministic aggregation later.
- Keep best-effort behavior: skip invalid entities without failing the batch.


