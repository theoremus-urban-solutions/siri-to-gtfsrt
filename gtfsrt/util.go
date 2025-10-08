package gtfsrt

import (
	"encoding/json"

	gtfs "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// UnmarshalPBFToProto parses PBF bytes into a protobuf FeedMessage.
func UnmarshalPBFToProto(b []byte) (*gtfs.FeedMessage, error) {
	pm := &gtfs.FeedMessage{}
	if err := proto.Unmarshal(b, pm); err != nil {
		return nil, err
	}
	return pm, nil
}

// ProtoToJSON returns canonical JSON for a protobuf FeedMessage suitable for comparing.
func ProtoToJSON(pb any) ([]byte, error) {
	m, ok := pb.(proto.Message)
	if !ok {
		return json.Marshal(pb)
	}
	marshaler := protojson.MarshalOptions{
		Multiline:       false,
		Indent:          "",
		UseProtoNames:   true,
		UseEnumNumbers:  true,
		EmitUnpopulated: false,
	}
	return marshaler.Marshal(m)
}
