package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"golang/convert"
	"golang/gtfsrt"
	sirixml "golang/siri/xml"
	"golang/types"
)

func main() {
	input := flag.String("input", "stdin", "file|url|stdin (url not yet supported)")
	path := flag.String("path", "", "PATH or URL when input is file or url")
	outfmt := flag.String("out", "gtfsrt-json", "gtfsrt-pbf|gtfsrt-json")
	kind := flag.String("type", "all", "trip-updates|vehicle-positions|alerts|all")
	output := flag.String("output", "", "output file or directory (stdout if empty)")
	split := flag.Bool("split", false, "when --type=all and output is a directory, write separate files")
	flag.Parse()

	_ = outfmt

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

	sd, err := sirixml.Decode(f)
	if err != nil {
		log.Fatalf("decode xml: %v", err)
	}

	entities, err := convert.ConvertSIRI(sd, types.DefaultOptions())
	if err != nil {
		log.Fatalf("convert: %v", err)
	}

	var msgs map[string]*gtfsrt.FeedMessage
	switch *kind {
	case "trip-updates":
		msgs = map[string]*gtfsrt.FeedMessage{"trip-updates": convert.BuildFeedMessage(filterByKind(entities, "trip_update"))}
	case "vehicle-positions":
		msgs = map[string]*gtfsrt.FeedMessage{"vehicle-positions": convert.BuildFeedMessage(filterByKind(entities, "vehicle_position"))}
	case "alerts":
		msgs = map[string]*gtfsrt.FeedMessage{"alerts": convert.BuildFeedMessage(filterByKind(entities, "alert"))}
	case "all":
		msgs = map[string]*gtfsrt.FeedMessage{
			"trip-updates":      convert.BuildFeedMessage(filterByKind(entities, "trip_update")),
			"vehicle-positions": convert.BuildFeedMessage(filterByKind(entities, "vehicle_position")),
			"alerts":            convert.BuildFeedMessage(filterByKind(entities, "alert")),
		}
	default:
		log.Fatalf("unsupported --type: %s", *kind)
	}

	switch *outfmt {
	case "gtfsrt-json":
		if *output == "" {
			first := true
			for name, m := range msgs {
				if !first {
					fmt.Println()
				}
				first = false
				fmt.Printf("# %s\n", name)
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				if err := enc.Encode(m); err != nil {
					log.Fatalf("encode json: %v", err)
				}
			}
		} else {
			fi, err := os.Stat(*output)
			if err == nil && fi.IsDir() {
				for name, m := range msgs {
					if !*split && *kind == "all" && name != "trip-updates" {
						continue
					}
					p := *output + "/" + name + ".json"
					if err := writeJSON(p, m); err != nil {
						log.Fatalf("write %s: %v", p, err)
					}
				}
			} else {
				if len(msgs) == 1 {
					for _, m := range msgs {
						if err := writeJSON(*output, m); err != nil {
							log.Fatalf("write %s: %v", *output, err)
						}
					}
				} else {
					out := map[string]*gtfsrt.FeedMessage(msgs)
					b, err := json.MarshalIndent(out, "", "  ")
					if err != nil {
						log.Fatalf("encode json: %v", err)
					}
					if err := os.WriteFile(*output, b, 0o644); err != nil {
						log.Fatalf("write %s: %v", *output, err)
					}
				}
			}
		}
	case "gtfsrt-pbf":
		log.Fatalf("PBF output not supported yet in placeholder build")
	default:
		log.Fatalf("unsupported --out: %s", *outfmt)
	}
}

func filterByKind(entities []types.Entity, kind string) []types.Entity {
	var out []types.Entity
	for _, e := range entities {
		if e.Kind == kind {
			out = append(out, e)
		}
	}
	return out
}

func writeJSON(path string, m *gtfsrt.FeedMessage) error {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
