Based on Entur's Kishar (embed @https://github.com/entur/kishar) — a SIRI → GTFS-Realtime converter reimplemented in Go with a flat package layout and a minimal CLI for converting SIRI ET/VM/SX XML into GTFS-RT-like JSON outputs.

What this is
- Lightweight Go library and CLI for transforming SIRI ServiceDelivery payloads:
  - ET (EstimatedTimetable) → TripUpdates
  - VM (VehicleMonitoring) → VehiclePositions
  - SX (SituationExchange) → Alerts

Usage (CLI)
1) Build
   - `go build -o kishar`
2) Run
   - `./kishar --input file --path siri.xml --type all --out gtfsrt-json --output outdir --split`
   - `cat siri.xml | ./kishar --type vehicle-positions --out gtfsrt-json`

Library API (import-less, flat package main)
- `ConvertSIRI(sd *ServiceDelivery, opts Options) ([]Entity, error)`
- `BuildFeedMessage(entities []Entity) *FeedMessage`
- `BuildPerDatasource(entities []Entity) map[string]*FeedMessage`

Library API details
- `ServiceDelivery`: In-memory representation of the SIRI envelope that can include any combination of the three delivery blocks:
  - `EstimatedTimetableDelivery` (ET) containing `EstimatedVehicleJourney` records
  - `VehicleMonitoringDelivery` (VM) containing `VehicleActivity` records
  - `SituationExchangeDelivery` (SX) containing `PtSituationElement` records
- `Options`:
  - `ETWhitelist`, `VMWhitelist`, `SXWhitelist`: optional datasource/participant filters (not wired in CLI yet)
  - `CloseToNextStopPercentage` and `CloseToNextStopDistance`: VM heuristics for proximity (reserved for future use)
  - `VMGracePeriod`: fallback TTL when message timestamps do not provide an explicit end time
- `Entity` (conversion output):
  - `ID`: stable key per entity (e.g., vehicleRef or tripId-startDate, situationNumber)
  - `Datasource`: origin/source if present in SIRI
  - `Kind`: one of `trip_update`, `vehicle_position`, `alert`
  - `Message`: GTFS-RT-like `FeedEntity` ready for aggregation/export
  - `TTL`: recommended time-to-live for caching/expiration
- `ConvertSIRI(...)`:
  - Iterates ET/VM/SX blocks in `ServiceDelivery` and maps each record to a corresponding `Entity`
  - Skips malformed/incomplete inputs without failing the whole batch
- `BuildFeedMessage(...)`:
  - Wraps a slice of `Entity` into a GTFS-RT-like `FeedMessage` (header + entity list)
- `BuildPerDatasource(...)`:
  - Groups `Entity` by `Datasource` and builds one `FeedMessage` per source (map)

What is “SIRI ServiceDelivery”?
- SIRI (Service Interface for Real-time Information) is a CEN standard for public-transport real-time data exchange.
- `ServiceDelivery` is the standard response/envelope element that a SIRI producer returns to a consumer. It may include zero or more of these delivery sections:
  - `EstimatedTimetableDelivery` (ET): predictions/updates for planned journeys and stop calls (drives TripUpdates)
  - `VehicleMonitoringDelivery` (VM): real-time vehicle positions and statuses (drives VehiclePositions)
  - `SituationExchangeDelivery` (SX): disruptions, messages, and advisories (drives Alerts)
- In XML, you’ll typically see `<Siri><ServiceDelivery>...` containing these nested blocks. This project decodes that payload and maps each block to the equivalent GTFS-RT concept.



