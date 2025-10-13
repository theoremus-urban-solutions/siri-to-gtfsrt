package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/theoremus-urban-solutions/siri-to-gtfsrt/converter"
	"github.com/theoremus-urban-solutions/siri-to-gtfsrt/formatter"
	"github.com/theoremus-urban-solutions/siri-to-gtfsrt/gtfsrt"
)

func main() {
	input := flag.String("input", "stdin", "file|stdin")
	path := flag.String("path", "", "PATH when input is file")
	outfmt := flag.String("out", "gtfsrt-json", "gtfsrt-pbf|gtfsrt-json")
	kind := flag.String("type", "all", "trip-updates|vehicle-positions|alerts|all")
	output := flag.String("output", "", "output file (stdout if empty)")
	flag.Parse()

	var f *os.File
	var err error
	switch *input {
	case "stdin":
		f = os.Stdin
	case "file":
		f, err = os.Open(*path)
		if err != nil {
			log.Fatalf("open: %v", err)
		}
		defer f.Close()
	default:
		log.Fatalf("unsupported input: %s", *input)
	}

	// Decode base64-encoded SIRI XML using streaming decoder
	sd, err := formatter.DecodeSIRIFromBase64(f)
	if err != nil {
		log.Fatalf("decode base64 xml: %v", err)
	}

	entities, err := converter.ConvertSIRI(sd, converter.DefaultOptions())
	if err != nil {
		log.Fatalf("convert: %v", err)
	}

	var msgs map[string]*gtfsrt.FeedMessage
	switch *kind {
	case "trip-updates":
		msgs = map[string]*gtfsrt.FeedMessage{"trip-updates": converter.BuildFeedMessage(filterByKind(entities, "trip_update"))}
	case "vehicle-positions":
		msgs = map[string]*gtfsrt.FeedMessage{"vehicle-positions": converter.BuildFeedMessage(filterByKind(entities, "vehicle_position"))}
	case "alerts":
		msgs = map[string]*gtfsrt.FeedMessage{"alerts": converter.BuildFeedMessage(filterByKind(entities, "alert"))}
	case "all":
		msgs = map[string]*gtfsrt.FeedMessage{
			"trip-updates":      converter.BuildFeedMessage(filterByKind(entities, "trip_update")),
			"vehicle-positions": converter.BuildFeedMessage(filterByKind(entities, "vehicle_position")),
			"alerts":            converter.BuildFeedMessage(filterByKind(entities, "alert")),
		}
	default:
		log.Fatalf("unsupported --type: %s", *kind)
	}

	switch *outfmt {
	case "gtfsrt-json":
		// Get the single message (should only be one based on type filter)
		var msg *gtfsrt.FeedMessage
		for _, m := range msgs {
			msg = m
			break
		}

		if *output == "" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(msg); err != nil {
				log.Fatalf("encode json: %v", err)
			}
		} else {
			b, err := json.MarshalIndent(msg, "", "  ")
			if err != nil {
				log.Fatalf("encode json: %v", err)
			}
			if err := os.WriteFile(*output, b, 0o644); err != nil {
				log.Fatalf("write %s: %v", *output, err)
			}
		}
	case "gtfsrt-pbf":
		if len(msgs) != 1 {
			log.Fatalf("when writing PBF, select a single --type")
		}
		for _, m := range msgs {
			b, err := gtfsrt.MarshalPBF(m)
			if err != nil {
				log.Fatalf("marshal pbf: %v", err)
			}
			if *output == "" {
				_, err = os.Stdout.Write(b)
			} else {
				err = os.WriteFile(*output, b, 0o644)
			}
			if err != nil {
				log.Fatalf("write: %v", err)
			}
		}
	default:
		log.Fatalf("unsupported --out: %s", *outfmt)
	}

	fmt.Fprintf(os.Stderr, "Converted %d entities\n", len(entities))
}

func filterByKind(entities []converter.Entity, kind string) []converter.Entity {
	var out []converter.Entity
	for _, e := range entities {
		if e.Kind == kind {
			out = append(out, e)
		}
	}
	return out
}
