### Project Overview (current scope)

Purpose
- Convert SIRI XML `ServiceDelivery` into GTFS-RT-like JSON using a Go library and a minimal CLI.

Flow
- Decode `<Siri><ServiceDelivery>...` XML
- Map ET → TripUpdates, VM → VehiclePositions, SX → Alerts
- Aggregate to `FeedMessage` per type (and optionally per datasource)

Out of scope (for now)
- No HTTP APIs, no Pub/Sub, no Redis, no metrics






