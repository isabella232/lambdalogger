package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/pkg/errors"
)

type rawEvent struct {
	AWSLogs struct {
		Data string
	}
}

type decodedEvent struct {
	// Owner is the account ID that generates this message
	Owner string

	// LogGroup is like /aws/lambda/logging-dev
	LogGroup string

	// MessageType Data messages will use the "DATA_MESSAGE" type.
	// Sometimes CloudWatch Logs may emit Kinesis records
	// with a "CONTROL_MESSAGE" type, mainly for checking if the
	// destination is reachable.
	MessageType string

	// LogEvents are the actual lines sent
	LogEvents []logEvent

	LogStream string
	// SubscriptionFilters []string // ignored
}

func (in *decodedEvent) FuncName() string {
	return path.Base(in.LogGroup)
}

type logEvent struct {
	ID        string
	Timestamp int
	Message   string
}

// decode decodes the base64/gzip compressed json data into a usable form
func (in *rawEvent) decode() (*decodedEvent, error) {
	compressed, err := base64.StdEncoding.DecodeString(in.AWSLogs.Data)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode raw input")
	}

	gread, err := gzip.NewReader(bytes.NewBuffer(compressed))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create GZIP reader")
	}

	out, err := ioutil.ReadAll(gread)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read in gzip data")
	}
	res := new(decodedEvent)
	if err := json.Unmarshal(out, res); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshall raw data")
	}

	return res, nil
}

// START RequestId: e13a2aca-f26e-11e8-990c-2df9ff70a259 Version: $LATEST
// END RequestId: e13a2aca-f26e-11e8-990c-2df9ff70a259
// REPORT RequestId: e13a2aca-f26e-11e8-990c-2df9ff70a259    Duration: 1.10 ms    Billed Duration: 100 ms     Memory Size: 256 MB    Max Memory Used: 26 MB
// func parseControlMesssage(raw string) *event {
// 	parts := strings.Fields(raw)
// 	if len(parts) < 3 {
// 		return nil
// 	}
// 	msgType := parts[0]
// 	evt := &event{
// 		Attributes: map[string]interface{
// 			"msg_type":   msgType,
// 			"request_id": parts[2],
// 		},
// 	}

// 	switch msgType {
// 	case "START":
// 		if len(parts) == 5 {
// 			evt.Attributes["version"] = parts[4]
// 		}
// 	case "END":
// 		// nothing
// 	case "REPORT":
// 		// 00: REPORT
// 		// 01: RequestId:
// 		// 02: e13a2aca-f26e-11e8-990c-2df9ff70a259
// 		// 03: Duration:
// 		// 04: 1.10
// 		// 05: ms
// 		// 06: Billed
// 		// 07: Duration:
// 		// 08: 100
// 		// 09: ms
// 		// 10: Memory
// 		// 11: Size:
// 		// 12: 256
// 		// 13: MB
// 		// 14: Max
// 		// 15: Memory
// 		// 16: Used:
// 		// 17: 26
// 		// 18: MB
// 		if len(parts) == 19 {
// 			evt.Attributes["dur"] = parts[4]
// 			evt.Attributes["dur_units"] = parts[5]
// 			evt.Attributes["billed_dur"] = parts[8]
// 			evt.Attributes["billed_dur_units"] = parts[9]
// 			evt.Attributes["mem_size"] = parts[12]
// 			evt.Attributes["mem_units"] = parts[13]
// 			evt.Attributes["max_mem"] = parts[17]
// 			evt.Attributes["max_mem_units"] = parts[18]
// 		}
// 	default:
// 		return nil
// 	}

// 	return evt
// }
