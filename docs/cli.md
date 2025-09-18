### CLI Specification — SIRI → GTFS-Realtime

Goal: A simple command-line tool wrapping the library to convert SIRI ET/VM/SX into GTFS-RT, reading from stdin/files/URLs and writing PBF or JSON.

### Commands
- `kishar` — Convert SIRI (XML) to GTFS-RT

### Flags
- Input selection:
  - `--input file|url|stdin` (default: stdin)
  - `--path PATH_OR_URL`
- Output selection:
  - `--out gtfsrt-pbf|gtfsrt-json` (default: gtfsrt-json)
  - `--output PATH` (default: stdout)
- Filters and options (placeholders for future wiring):
  - `--whitelist-et csv`
  - `--whitelist-vm csv`
  - `--whitelist-sx csv`
  - `--vm-close-percentage int` (default: 95)
  - `--vm-close-distance int` meters (default: 500)
  - `--vm-grace-period duration` (default: 5m)

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
go build -o kishar

# Convert SIRI XML from stdin to GTFS-RT trip-updates in JSON
cat siri.xml | ./kishar --type trip-updates --out gtfsrt-json > trip-updates.json

# Convert SIRI ET/VM/SX XML file to three GTFS-RT JSON files
./kishar --input file --path siri.xml --type all --out gtfsrt-json --output outdir --split

# Convert and print Vehicle Positions JSON
./kishar --input file --path siri.xml --type vehicle-positions --out gtfsrt-json | jq

# Apply datasource whitelists (placeholders)
./kishar --input file --path siri.xml --whitelist-et RUT,BNR --whitelist-vm RUT --whitelist-sx ENT
```

### Exit codes
- 0 on success with any produced output
- 1 on invalid input/flags or decode failures
- 2 when input is valid but yields zero entities and `--strict` is set

### Notes
- PBF output is not yet implemented in the flat placeholder build (`--out gtfsrt-pbf` will exit with an error).
- For streaming use later, this CLI will be used as a building block in services; it should be stateless and fast.




