# siri-to-gtfsrt

A Go library and CLI tool for converting SIRI (Service Interface for Real Time Information) transit data to GTFS-Realtime format.

Based on Entur's Kishar (https://github.com/entur/kishar) — reimplemented in Go with a modular package structure.

**Sister project**: [gtfsrt-to-siri](https://github.com/theoremus-urban-solutions/gtfsrt-to-siri) - Convert GTFS-RT to SIRI (inverse operation)

## Features

- **Comprehensive**: Supports Vehicle Monitoring (VM), Estimated Timetable (ET), and Situation Exchange (SX)
- **Library-First**: Clean, modular API designed for server integration
- **CLI Tools**: Command-line tools for standalone conversion and testing
- **Well-Structured**: Modular package design with clear separation of concerns

## Installation

```bash
go get github.com/theoremus-urban-solutions/siri-to-gtfsrt
```

## Quick Start

### As a Library

```go
package main

import (
    "os"
    "log"
    
    "github.com/theoremus-urban-solutions/siri-to-gtfsrt/converter"
    "github.com/theoremus-urban-solutions/siri-to-gtfsrt/formatter"
    "github.com/theoremus-urban-solutions/siri-to-gtfsrt/gtfsrt"
)

func main() {
    // 1. Parse SIRI XML
    file, _ := os.Open("siri-vm.xml")
    defer file.Close()
    
    serviceDelivery, err := formatter.DecodeSIRI(file)
    if err != nil {
        log.Fatal(err)
    }
    
    // 2. Convert to GTFS-RT
    entities, err := converter.ConvertSIRI(serviceDelivery, converter.DefaultOptions())
    if err != nil {
        log.Fatal(err)
    }
    
    // 3. Build feed message
    feedMessage := converter.BuildFeedMessage(entities)
    
    // 4. Serialize to protobuf
    pbfBytes, err := gtfsrt.MarshalPBF(feedMessage)
    if err != nil {
        log.Fatal(err)
    }
    
    // Write to file or serve via HTTP
    os.WriteFile("gtfsrt.pb", pbfBytes, 0644)
}
```

### As a CLI Tool

Install the CLI:

```bash
go install github.com/theoremus-urban-solutions/siri-to-gtfsrt/cmd/siri-to-gtfsrt@latest
```

Convert SIRI XML to GTFS-RT:

```bash
# Convert to JSON
siri-to-gtfsrt --input=file --path=siri-vm.xml --type=vehicle-positions --out=gtfsrt-json

# Convert to protobuf
siri-to-gtfsrt --input=file --path=siri-vm.xml --type=vehicle-positions --out=gtfsrt-pbf > output.pb

# Process from stdin
cat siri.xml | siri-to-gtfsrt --type=vehicle-positions --out=gtfsrt-pbf > output.pb
```

## Architecture

The library follows a clean, modular architecture:

```
┌─────────────────────────────────────┐
│         Your Application            │
└──────────┬──────────────────────────┘
           │
           ├──► formatter.DecodeSIRI()     (Parse SIRI XML)
           │      └─► siri.ServiceDelivery
           │
           ├──► converter.ConvertSIRI()    (Business logic)
           │      └─► []converter.Entity
           │
           ├──► converter.BuildFeedMessage() (Build feed)
           │      └─► gtfsrt.FeedMessage
           │
           └──► gtfsrt.MarshalPBF()        (Serialize)
                  └─► []byte
```

### Package Overview

- **`siri/`**: SIRI domain types (ServiceDelivery, VM, ET, SX)
- **`gtfsrt/`**: GTFS-RT types and protobuf operations
- **`converter/`**: Conversion business logic
- **`formatter/`**: Input/output formatting (XML, JSON)
- **`cmd/`**: CLI applications

## Usage Examples

### Filtering by Datasource

```go
// Convert and organize by datasource
entities, _ := converter.ConvertSIRI(sd, converter.DefaultOptions())
byDatasource := converter.BuildPerDatasource(entities)

// Access specific datasource
operatorAFeed := byDatasource["OPERATOR_A"]
operatorBFeed := byDatasource["OPERATOR_B"]
```

### Custom Options

```go
opts := converter.Options{
    VMGracePeriod: 10 * time.Minute,
    ETWhitelist: []string{"ROUTE_1", "ROUTE_2"},
    VMWhitelist: []string{"VEHICLE_123"},
    CloseToNextStopPercentage: 90,
    CloseToNextStopDistance: 300,
}

entities, _ := converter.ConvertSIRI(sd, opts)
```

## CLI Reference

### siri-to-gtfsrt

Convert SIRI XML to GTFS-Realtime.

**Flags:**

- `--input`: Input source (`file`, `stdin`) [default: `stdin`]
- `--path`: Path to input file (when `--input=file`)
- `--type`: Entity type (`trip-updates`, `vehicle-positions`, `alerts`, `all`) [default: `all`]
- `--out`: Output format (`gtfsrt-json`, `gtfsrt-pbf`) [default: `gtfsrt-pbf`]
- `--output`: Output file or directory [default: stdout]
- `--split`: Write separate files when `--type=all` and output is a directory

**Examples:**

```bash
# Single feed type to stdout
cat vm.xml | siri-to-gtfsrt --type=vehicle-positions > vp.pb

# All types to separate files
siri-to-gtfsrt --input=file --path=siri.xml --type=all --output=./out --split

# JSON for debugging
siri-to-gtfsrt --input=file --path=vm.xml --type=vehicle-positions --out=gtfsrt-json | jq .
```

### gtfsrt-diff

Compare two GTFS-RT feeds for regression testing.

**Flags:**

- `--ours`: Path to output feed (file or directory)
- `--golden`: Path to golden/reference feed (file or directory)
- `--format`: Feed format (`json`, `pbf`) [default: `json`]
- `--html`: Output HTML report path

**Example:**

```bash
gtfsrt-diff --ours=out/ --golden=testdata/golden/ --format=json --html=report.html
```

## Building from Source

```bash
# Clone the repository
git clone https://github.com/theoremus-urban-solutions/siri-to-gtfsrt
cd siri-to-gtfsrt

# Build CLI tools
go build -o bin/siri-to-gtfsrt ./cmd/siri-to-gtfsrt
go build -o bin/gtfsrt-diff ./cmd/gtfsrt-diff

# Run tests
go test ./...
```

## What is SIRI?

SIRI (Service Interface for Real-time Information) is a CEN standard for public transport real-time data exchange. The `ServiceDelivery` element is the standard response envelope that may include:

- **EstimatedTimetableDelivery (ET)**: Predictions/updates for planned journeys and stop calls → TripUpdates
- **VehicleMonitoringDelivery (VM)**: Real-time vehicle positions and statuses → VehiclePositions
- **SituationExchangeDelivery (SX)**: Disruptions, messages, and advisories → Alerts

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

[Your License Here]

## Related Projects

- [gtfsrt-to-siri](https://github.com/theoremus-urban-solutions/gtfsrt-to-siri) - Convert GTFS-RT to SIRI (inverse operation)
- [MobilityData GTFS-RT Bindings](https://github.com/MobilityData/gtfs-realtime-bindings) - Official GTFS-RT protobuf bindings
- [Entur Kishar](https://github.com/entur/kishar) - Original Java implementation