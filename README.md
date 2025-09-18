Based on Entur's Kishar (embed @https://github.com/entur/kishar) — a SIRI → GTFS-Realtime converter reimplemented in Go with a flat package layout and a minimal CLI for converting SIRI ET/VM/SX XML into GTFS-RT-like JSON outputs.

What this is
- Lightweight Go library and CLI for transforming SIRI ServiceDelivery payloads:
  - ET (EstimatedTimetable) → TripUpdates
  - VM (VehicleMonitoring) → VehiclePositions
  - SX (SituationExchange) → Alerts

Usage (CLI)
1) Build
   - `go build -o siri-to-gtfs ./cmd/siri-to-gtfs`
2) Run
   - `./siri-to-gtfs --input file --path siri.xml --type all --out gtfsrt-json --output outdir --split`
   - `cat siri.xml | ./siri-to-gtfs --type vehicle-positions --out gtfsrt-json`

Library API (import-less, flat package main)
- `ConvertSIRI(sd *ServiceDelivery, opts Options) ([]Entity, error)`
- `BuildFeedMessage(entities []Entity) *FeedMessage`
- `BuildPerDatasource(entities []Entity) map[string]*FeedMessage`

Library API details
- `ServiceDelivery`: In-memory representation of the SIRI envelope that can include any combination of the three delivery blocks:
  - `EstimatedTimetableDelivery` (ET) containing `EstimatedVehicleJourney` records
  - `VehicleMonitoringDelivery` (VM) containing `VehicleActivity` records
  - `SituationExchangeDelivery` (SX) containing `PtSituationElement` records
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

### CLI Specification — SIRI → GTFS-Realtime

Goal: A simple command-line tool wrapping the library to convert SIRI ET/VM/SX into GTFS-RT, reading from stdin/files/URLs and writing PBF or JSON.

### Commands
- `siri-to-gtfs` — Convert SIRI (XML) to GTFS-RT

### Flags
- Input selection:
  - `--input file|url|stdin` (default: stdin)
  - `--path PATH_OR_URL`
- Output selection:
  - `--out gtfsrt-pbf|gtfsrt-json` (default: gtfsrt-json)
  - `--output PATH` (default: stdout)

### Behavior
1) Read SIRI payload (XML); decode to internal SIRI structs.
2) Run `convert.ConvertSIRI` with flags → entities.
3) Aggregate to one GTFS-RT `FeedMessage` per type or combined:
   - Default combined single `FeedMessage` that contains all entity types is not standard. Instead, provide explicit output modes:
     - `--type trip-updates|vehicle-positions|alerts|all` (default: all)
     - For `all`, emit multiple outputs when writing to files: `trip-updates.*`, `vehicle-positions.*`, `alerts.*`. For stdout, emit one after another separated by boundaries or use `--split` to write directory outputs.
4) Write output as PBF or JSON.

### Examples
```bash
# Build the CLI
go build -o siri-to-gtfs ./cmd/siri-to-gtfs

# Convert SIRI XML from stdin to GTFS-RT trip-updates in JSON
cat siri.xml | ./siri-to-gtfs --type trip-updates --out gtfsrt-json > trip-updates.json

# Convert SIRI ET/VM/SX XML file to three GTFS-RT JSON files
./siri-to-gtfs --input file --path siri.xml --type all --out gtfsrt-json --output outdir --split

# Convert and print Vehicle Positions JSON
./siri-to-gtfs --input file --path siri.xml --type vehicle-positions --out gtfsrt-json | jq
```

### Exit codes
- 0 on success with any produced output
- 1 on invalid input/flags or decode failures
- 2 when input is valid but yields zero entities and `--strict` is set

### Notes
- PBF output is not yet implemented in the flat placeholder build (`--out gtfsrt-pbf` will exit with an error).
- For streaming use later, this CLI will be used as a building block in services; it should be stateless and fast.






