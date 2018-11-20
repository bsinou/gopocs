package bleve

import (
	"fmt"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/rs/xid"
)

// SimpleTestServer is a canonical test implementation of the Log server
type SimpleTestServer struct {
	Index bleve.Index
	idgen xid.ID
}

type indexableLog struct {
	LogMessage
}

func (l *indexableLog) String() string {
	return fmt.Sprintf("%d - %s - %s - %s - %s - %s\n", l.Ts, l.Logger, l.Level, l.UserName, l.MsgId, l.Msg)
}

// LogMessage is a basic log format used to transmit log messages to clients.
type LogMessage struct {
	Ts            int32  `protobuf:"varint,1,opt,name=Ts" json:"Ts,omitempty"`
	Level         string `protobuf:"bytes,2,opt,name=Level" json:"Level,omitempty"`
	Logger        string `protobuf:"bytes,3,opt,name=Logger" json:"Logger,omitempty"`
	Msg           string `protobuf:"bytes,4,opt,name=Msg" json:"Msg,omitempty"`
	MsgId         string `protobuf:"bytes,5,opt,name=MsgId" json:"MsgId,omitempty"`
	UserName      string `protobuf:"bytes,6,opt,name=UserName" json:"UserName,omitempty"`
	RemoteAddress string `protobuf:"bytes,10,opt,name=RemoteAddress" json:"RemoteAddress,omitempty"`
}

// NewSimpleTestServer creates and configures a new Bleve instance if none has already been configured.
func NewSimpleTestServer(bleveIndexPath string, deleteOnClose ...bool) (*SimpleTestServer, error) {

	index, err := bleve.Open(bleveIndexPath)
	if err == nil {
		// Already existing, no needs to create
		return &SimpleTestServer{Index: index}, nil
	}

	// Creates the new index and initialises the server
	mapping := bleve.NewIndexMapping()
	defaultMapping := bleve.NewDocumentMapping()

	// Specific fields
	standardFieldMapping := bleve.NewTextFieldMapping()
	defaultMapping.AddFieldMappingsAt("level", standardFieldMapping)

	dateFieldMapping := bleve.NewDateTimeFieldMapping()
	defaultMapping.AddFieldMappingsAt("ts", dateFieldMapping)

	msgFieldMapping := bleve.NewTextFieldMapping()
	defaultMapping.AddFieldMappingsAt("msg", msgFieldMapping)

	mapping.AddDocumentMapping("simpletest", defaultMapping)

	if bleveIndexPath == "" {
		index, err = bleve.NewMemOnly(mapping)
	} else {
		index, err = bleve.New(bleveIndexPath, mapping)
	}
	return &SimpleTestServer{Index: index}, nil
}

// PutLog stores a new log msg. It expects a map[string]string.
func (s *SimpleTestServer) PutLog(line map[string]string) error {
	msg, err := toIndexableMsg(line)
	if err != nil {
		return err
	}

	err = s.Index.Index(xid.New().String(), msg)
	if err != nil {
		return err
	}
	return nil
}

func toIndexableMsg(line map[string]string) (*indexableLog, error) {

	// fmt.Println("\n\n## [DEBUG] ## Marshall Log line B4 indexing:")
	// for k, v := range line {
	// 	fmt.Printf("## [DEBUG] [%s] - [%s]\n", k, v)
	// }

	msg := &indexableLog{}

	for k, val := range line {
		switch k {
		case "ts":
			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				return nil, err
			}
			msg.Ts = convertTimeToTs(t)
		case "level":
			msg.Level = val
		case "MsgId":
			msg.MsgId = val
		case "logger": // name of the service that is currently logging.
			msg.Logger = val
		case "UserName":
			msg.UserName = val
		default:
			break
		}
	}

	// Concatenate msg and error in the full text msg field.
	text := ""
	if m, ok := line["msg"]; ok {
		text = m
	}
	if m, ok := line["error"]; ok {
		text += " - " + m
	}
	msg.Msg = text

	// fmt.Printf("## Msg created: %s\n", msg.String())

	return msg, nil
}

// Single entry point to convert time.Time to Unix timestamps defined as int32
func convertTimeToTs(ts time.Time) int32 {
	return int32(ts.Unix())
}

// Single entry point to convert Unix timestamps defined as int32 to time.Time
func convertTsToTime(ts int32) time.Time {
	return time.Unix(int64(ts), 0)
}