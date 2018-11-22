package bleve

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/blevesearch/bleve/index/store/boltdb"
	"github.com/blevesearch/bleve/index/upsidedown"
	"github.com/rs/xid"
)

// TestServer is a generic index server
type TestServer struct {
	name  string
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

// NewDefaultServer creates and configures a new Default Bleve instance if none has already been configured.
func NewDefaultServer(bleveIndexPath, sname string) (*TestServer, error) {

	index, err := bleve.Open(bleveIndexPath)
	if err == nil {
		// Already existing, no needs to create
		return &TestServer{Index: index, name: sname}, nil
	}

	// Creates the new index and initialises the server
	mapping := bleve.NewIndexMapping()
	defaultMapping := bleve.NewDocumentMapping()

	mapping.AddDocumentMapping("simpletest", defaultMapping)

	if bleveIndexPath == "" {
		index, err = bleve.NewMemOnly(mapping)
	} else {
		index, err = bleve.New(bleveIndexPath, mapping)
	}
	return &TestServer{Index: index, name: sname}, nil
}

// NewDefaultServerWithMapping creates and configures a new Default Bleve instance if none has already been configured.
func NewDefaultServerWithMapping(bleveIndexPath, sname string) (*TestServer, error) {

	index, err := bleve.Open(bleveIndexPath)
	if err == nil {
		// Already existing, no needs to create
		return &TestServer{Index: index, name: sname}, nil
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

	mapping.AddDocumentMapping("explicitmapping", defaultMapping)

	if bleveIndexPath == "" {
		index, err = bleve.NewMemOnly(mapping)
	} else {
		index, err = bleve.New(bleveIndexPath, mapping)
	}
	return &TestServer{Index: index, name: sname}, nil
}

// NewServerOnBolt creates and configures a bleve instance that use BoltDB if none has already been configured.
func NewServerOnBolt(bleveIndexPath, sname string) (*TestServer, error) {

	index, err := bleve.Open(bleveIndexPath)
	if err == nil {
		// Already existing, no needs to create
		return &TestServer{Index: index, name: sname}, nil
	}

	// Creates the new index and initialises the server
	mapping := bleve.NewIndexMapping()
	defaultMapping := bleve.NewDocumentMapping()

	// // Specific fields
	// standardFieldMapping := bleve.NewTextFieldMapping()
	// defaultMapping.AddFieldMappingsAt("level", standardFieldMapping)
	// dateFieldMapping := bleve.NewDateTimeFieldMapping()
	// defaultMapping.AddFieldMappingsAt("ts", dateFieldMapping)
	// msgFieldMapping := bleve.NewTextFieldMapping()
	// defaultMapping.AddFieldMappingsAt("msg", msgFieldMapping)

	mapping.AddDocumentMapping("dft-on-bolt", defaultMapping)
	index, err = bleve.NewUsing(bleveIndexPath, mapping, upsidedown.Name, boltdb.Name, nil)

	// config := map[string]interface{}{
	//     "nosync": true,
	// }
	// index, err := bleve.NewUsing(bleveIndexPath, mapping, upsidedown.Name, boltdb.Name, config)

	return &TestServer{Index: index, name: sname}, nil
}

// NewServerOnScorch creates and configures a bleve instance that use BoltDB if none has already been configured.
func NewServerOnScorch(bleveIndexPath, sname string) (*TestServer, error) {

	index, err := bleve.Open(bleveIndexPath)
	if err == nil {
		// Already existing, no needs to create
		return &TestServer{Index: index, name: sname}, nil
	}

	// Creates the new index and initialises the server
	mapping := bleve.NewIndexMapping()
	defaultMapping := bleve.NewDocumentMapping()

	mapping.AddDocumentMapping("dft-on-scorch", defaultMapping)
	index, err = bleve.NewUsing(bleveIndexPath, mapping, scorch.Name, boltdb.Name, nil)

	return &TestServer{Index: index, name: sname}, nil
}

// TestServer specific methods.

// GetName simply return a name for this server.
func (s *TestServer) GetName() string {
	return s.name
}

// PutLog stores a new log msg. It expects a map[string]string.
func (s *TestServer) PutLog(line map[string]string) error {
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

func Json2map(line string) map[string]string {
	var rawdata map[string]interface{}
	// var data map[string]string
	err := json.Unmarshal([]byte(line), &rawdata)
	if err != nil {
		fmt.Println(err)
	}

	data := make(map[string]string)
	for k, v := range rawdata {
		switch v := v.(type) {
		case string:
			data[k] = v
		default:
			fmt.Printf("Cannot unmarshall object for key %s\n", k)
		}
	}

	return data
}
