package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type humioConfig struct {
	Token      string
	Repository string
	Parser     string
	Endpoint   string `default:"https://cloud.humio.com"`
}

/*
https://docs.humio.com/api/ingest-api/
  {
    "type": "accesslog",
    "fields": {
      "host": "webhost1"
    },
    "messages": [
       "192.168.1.21 - user1 [02/Nov/2017:13:48:26 +0000] \"POST /humio/api/v1/dataspaces/humio/ingest HTTP/1.1\" 200 0 \"-\" \"useragent\" 0.015 664 0.015",
       "192.168.1.49 - user1 [02/Nov/2017:13:48:33 +0000] \"POST /humio/api/v1/dataspaces/developer/ingest HTTP/1.1\" 200 0 \"-\" \"useragent\" 0.014 657 0.014",
       "192.168.1..21 - user2 [02/Nov/2017:13:49:09 +0000] \"POST /humio/api/v1/dataspaces/humio HTTP/1.1\" 200 0 \"-\" \"useragent\" 0.013 565 0.013",
       "192.168.1.54 - user1 [02/Nov/2017:13:49:10 +0000] \"POST /humio/api/v1/dataspaces/humio/queryjobs HTTP/1.1\" 200 0 \"-\" \"useragent\" 0.015 650 0.015"
    ]
  }
*/
type humioMsg struct {
	// Type The parser Humio will use to parse the messages
	Type string `json:"type,omitempty"`
	// Fields Annotate each of the messages with these key-values. Values must be strings.
	Fields map[string]string `json:"fields,omitempty"`
	// Tags Annotate each of the messages with these key-values as Tags. Please see other documentation on tags before using.
	Tags map[string]interface{} `json:"tags,omitempty"`

	// Messages	The raw strings representing the events. Each string will be parsed by the parser specified by type.
	Messages []string `json:"messages,omitempty"`
}

type event struct {
	RawString  string                 `json:"raw_string,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

func newHumioMsg(decoded *decodedEvent) *humioMsg {
	msg := &humioMsg{
		Tags: map[string]interface{}{
			"aws_account_id": decoded.Owner,
			"function":       decoded.FuncName(),
		},
		Fields: map[string]string{
			"log_stream": decoded.LogStream,
			"log_group":  decoded.LogGroup,
		},
	}

	for k, v := range config.ExtraTags {
		msg.Tags[k] = v
	}

	for k, v := range config.ExtraFields {
		msg.Fields[k] = v
	}

	if config.Humio.Parser != "" {
		msg.Type = config.Humio.Parser
	} else {
		msg.Type = decoded.FuncName()
	}

	for _, le := range decoded.LogEvents {
		msg.Messages = append(msg.Messages, le.Message)
	}
	return msg
}

func sendStrings(out *humioMsg) (int, error) {
	url := fmt.Sprintf("%s/api/v1/dataspaces/%s/ingest-messages", config.Humio.Endpoint, config.Humio.Repository)
	return send(url, &[]*humioMsg{out})
}

func send(url string, input interface{}) (int, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to marshal json")
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return 0, errors.Wrap(err, "Failed to make a new request object")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Humio.Token))
	rsp, err := client.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to make the request")
	}
	if rsp.Body != nil {
		defer rsp.Body.Close()
	}

	if rsp.StatusCode != http.StatusOK {
		val, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return rsp.StatusCode, errors.Wrap(err, "Failed to read response body")
		}
		return rsp.StatusCode, errors.New(string(val))
	}

	return rsp.StatusCode, nil
}

func sendEvents(events []event) (int, error) {
	out := []struct {
		Tags   map[string]string `json:"tags"`
		Events []event           `json:"events"`
	}{{
		Tags: map[string]string{
			"function": "lambdalogger",
		},
		Events: events,
	}}
	url := fmt.Sprintf("%s/api/v1/dataspaces/%s/ingest", config.Humio.Endpoint, config.Humio.Repository)
	return send(url, &out)
}
