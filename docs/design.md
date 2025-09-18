### Go Library Design — SIRI (ET/VM/SX) → GTFS-Realtime

Scope
- Library + CLI only. No HTTP server, no Pub/Sub, no Redis. Focus on decoding SIRI XML and mapping to GTFS-RT-like JSON structures.

Inputs and outputs
- Input: SIRI XML `ServiceDelivery` containing any of ET/VM/SX blocks
- Output: In-memory GTFS-RT-like `FeedMessage` and `FeedEntity` objects; JSON rendering supported by the CLI

Library API (flat)
- `ConvertSIRI(sd *ServiceDelivery, opts Options) ([]Entity, error)`
- `BuildFeedMessage(entities []Entity) *FeedMessage`
- `BuildPerDatasource(entities []Entity) map[string]*FeedMessage`

Entity rules
- IDs: ET → `tripId[-startDate]`; VM → `vehicleRef` or `tripId-startDate`; SX → `situationNumber`
- TTL: ET → latest of time fields; VM → `validUntilTime` or `VMGracePeriod`; SX → validity end or 365d

Error handling
- Best-effort conversion; skip malformed records without failing the batch

Future work (optional)
- Protobuf bindings for GTFS-RT PBF
- Whitelist filtering in CLI flags