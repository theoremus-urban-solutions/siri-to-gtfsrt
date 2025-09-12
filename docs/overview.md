### Kishar — System Overview (for Go Rewrite)

**Purpose**: Ingest SIRI real-time data from Google Pub/Sub, transform to GTFS-Realtime (GTFS-RT), publish binary/JSON feeds over HTTP, and expose health and Prometheus metrics.

### High-level flow
- **Ingress (Pub/Sub)**: Subscribes to three topics (configurable):
  - **SIRI-ET** (Estimated Timetable → TripUpdates)
  - **SIRI-VM** (Vehicle Monitoring → VehiclePositions)
  - **SIRI-SX** (Situation Exchange → Alerts)
- **Parse/Map**: Convert Avro SIRI payloads to GTFS-RT entities:
  - `TripUpdate`, `VehiclePosition`, `Alert`
  - Applies datasource whitelists per feed type
  - Derives trip/vehicle descriptors and stop-time updates
  - Computes TTLs from SIRI timestamps; defaults used when missing
- **Store**: Upsert serialized GTFS-RT `FeedEntity` bytes into Redis (keyed by `CompositeKey{id,datasource}`) with per-entity TTL. In non-Redis mode, an in-memory map is used.
- **Aggregate/Publish**: On a 10s timer, read Redis and build full GTFS-RT `FeedMessage`s (total and per-datasource variants) for:
  - TripUpdates, VehiclePositions, Alerts
- **Serve**: HTTP endpoints return current GTFS-RT feeds (PBF by default, JSON if `Content-Type: application/json`).
- **Observe**: Health endpoints and Prometheus metrics.

### External interfaces
- HTTP (Jetty on `kishar.incoming.port`, default 8888):
  - `/api/trip-updates`
  - `/api/vehicle-positions`
  - `/api/alerts`
    - Optional header `datasource` filters to per-datasource feed
    - `Content-Type: application/json` returns JSON; otherwise PBF bytes
  - `/internal/debug/status` (entity counts)
  - `/internal/debug/reset` (clears Redis)
- Health & metrics:
  - `/health/ready`, `/health/up`, `/health/healthy`
  - `/health/scrape` (Prometheus exposition)
- Ingress (Google Pub/Sub via Camel URIs):
  - `kishar.pubsub.topic.et`, `kishar.pubsub.topic.vm`, `kishar.pubsub.topic.sx`

### Core components (Java reference)
- `PubSubRoute`: Wires Pub/Sub → parse → register flows and metrics.
- `GtfsRtProviderRoute`: Exposes REST endpoints and periodic aggregation (10s) to build `FeedMessage`s.
- `SiriToGtfsRealtimeService`: Central transformer/registry:
  - Converts SIRI ET/VM/SX objects into GTFS-RT builders
  - Applies whitelists; registers metrics
  - Computes entity TTLs and writes to Redis via `RedisService`
  - Aggregates and serves `FeedMessage`s (total/per-datasource)
- `GtfsRtMapper`: Field-level mapping from SIRI Avro → GTFS-RT, including:
  - Trip/Vehicle descriptors; stop-time updates with delay calc
  - Vehicle position including status, occupancy, congestion, odometer
  - Date/time formatting (GTFS-RT startDate/startTime)
- `AlertFactory`: Translates SIRI SX situations to GTFS-RT alerts (informed entities for trips/routes/stops, validity, URLs).
- `RedisService`: Writes/reads serialized `FeedEntity` bytes with TTL (Redisson) or in-memory fallback.
- `PrometheusMetricsService`: Custom counters and gauges for inbound requests/entities and GTFS-RT totals.
- GraphQL helpers (`JourneyPlannerGraphQLClient`, `ServiceJourneyService`):
  - Optional lookup of `ServiceJourney` from `DatedServiceJourney` to resolve trip IDs and dates.

### Configuration (key properties)
- Pub/Sub:
  - `spring.cloud.gcp.pubsub.project-id`
  - `spring.cloud.gcp.pubsub.credentials.location`
  - `kishar.pubsub.enabled=true|false`
  - `kishar.pubsub.topic.et|vm|sx`
- Server:
  - `kishar.incoming.port` (default 8888)
- Data whitelists:
  - `kishar.datasource.et.whitelist`
  - `kishar.datasource.vm.whitelist`
  - `kishar.datasource.sx.whitelist`
- Vehicle proximity heuristics:
  - `kishar.settings.vm.close.to.stop.percentage`
  - `kishar.settings.vm.close.to.stop.distance`
- Redis:
  - `kishar.redis.enabled`
  - `kishar.redis.host`, `kishar.redis.port`, `kishar.redis.password`
- Journey Planner GraphQL (optional):
  - `kishar.journeyplanner.url`, `kishar.journeyplanner.EtClientName`

### Data model summary
- Keys: `CompositeKey{id,datasource}` serialized to JSON for Redis keys.
- Values: Serialized GTFS-RT `FeedEntity` bytes plus TTL (`google.protobuf.Duration`).
- Aggregation builds `FeedMessage` per type, and per datasource map.

### Go rewrite guidance (scope outline)
- Replace Spring Boot + Camel with Go:
  - Pub/Sub subscriber clients for ET/VM/SX
  - HTTP server for endpoints and Prometheus `/metrics`
  - Timer to aggregate feeds every 10s
- Re-implement mappers and factories:
  - SIRI Avro inputs → consider equivalent Go structs or JSON decoding
  - GTFS-RT protobuf via `google.golang.org/protobuf`
  - Alert, TripUpdate, VehiclePosition assembly mirrors Java logic
- Storage:
  - Redis client (e.g., `go-redis`), TTL per key; in-memory fallback for local/dev
- Config:
  - Env vars or `yaml` with parity to current properties
- Observability:
  - Prometheus counters/gauges compatible with current metrics

This document describes current behavior to preserve functional parity during the Go migration.



