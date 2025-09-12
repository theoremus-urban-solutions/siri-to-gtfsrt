### CLI Specification — SIRI → GTFS-Realtime

Goal: A simple command-line tool wrapping the library to convert SIRI ET/VM/SX into GTFS-RT, reading from stdin/files/URLs and writing PBF or JSON.

### Commands
- `kishar-convert siri2gtfs` — Convert SIRI to GTFS-RT

### Flags
- Input selection:
  - `--input file|url|stdin` (default: stdin)
  - `--path PATH_OR_URL`
  - `--format json|xml` (default: json)
- Output selection:
  - `--out gtfsrt-pbf|gtfsrt-json` (default: gtfsrt-pbf)
  - `--output PATH` (default: stdout)
- Filters and options:
  - `--whitelist-et csv`
  - `--whitelist-vm csv`
  - `--whitelist-sx csv`
  - `--vm-close-percentage int` (default: 95)
  - `--vm-close-distance int` meters (default: 500)
  - `--vm-grace-period duration` (default: 5m)

### Behavior
1) Read SIRI payload (JSON initially); decode to internal SIRI structs.
2) Run `convert.ConvertSIRI` with flags → entities.
3) Aggregate to one GTFS-RT `FeedMessage` per type or combined:
   - Default combined single `FeedMessage` that contains all entity types is not standard. Instead, provide explicit output modes:
     - `--type trip-updates|vehicle-positions|alerts|all` (default: all)
     - For `all`, emit multiple outputs when writing to files: `trip-updates.*`, `vehicle-positions.*`, `alerts.*`. For stdout, emit one after another separated by boundaries or use `--split` to write directory outputs.
4) Write output as PBF or JSON.

### Examples
```bash
# Convert SIRI JSON from stdin to GTFS-RT trip-updates in PBF
cat siri.json | kishar-convert siri2gtfs --type trip-updates --out gtfsrt-pbf > trip-updates.pbf

# Convert SIRI ET/VM/SX JSON file to three GTFS-RT files
kishar-convert siri2gtfs --input file --path siri.json --type all --out gtfsrt-pbf --output outdir --split

# Convert and print JSON
kishar-convert siri2gtfs --input file --path siri.json --type vehicle-positions --out gtfsrt-json | jq

# Apply datasource whitelists
kishar-convert siri2gtfs --input file --path siri.json --whitelist-et RUT,BNR --whitelist-vm RUT --whitelist-sx ENT
```

### Exit codes
- 0 on success with any produced output
- 1 on invalid input/flags or decode failures
- 2 when input is valid but yields zero entities and `--strict` is set

### Notes
- XML support can be added later with `--format xml`.
- For streaming use later, this CLI will be used as a building block in services; it should be stateless and fast.


