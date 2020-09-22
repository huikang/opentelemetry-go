package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	apitrace "go.opentelemetry.io/otel/api/trace"
	coltracepb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/collector/trace/v1"
	"go.opentelemetry.io/otel/exporters/otlp/internal/transform"

	"go.opentelemetry.io/otel/label"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"

	_ "go.opentelemetry.io/otel/exporters/otlp/protocol/http"
)

var payloadType = "protostr" // "raw"

func makePayload() string {
	str := ""
	// March 31, 2020 5:01:26 1234nanos (UTC)
	startTime := time.Unix(1585674086, 1234)
	endTime := startTime.Add(10 * time.Second)
	spanData := &export.SpanData{
		SpanContext: apitrace.SpanContext{
			TraceID: apitrace.ID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
			SpanID:  apitrace.SpanID{0x00, 0x01, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
		},
		SpanKind:     apitrace.SpanKindServer,
		ParentSpanID: apitrace.SpanID{0xEF, 0xEE, 0xED, 0xEC, 0xEB, 0xEA, 0xE9, 0xE8},
		Name:         "span data to span data",
		StartTime:    startTime,
		EndTime:      endTime,
		MessageEvents: []export.Event{
			{Time: startTime,
				Attributes: []label.KeyValue{
					label.Uint64("CompressedByteSize", 512),
				},
			},
			{Time: endTime,
				Attributes: []label.KeyValue{
					label.String("MessageEventType", "Recv"),
				},
			},
		},
		/*
			Links: []apitrace.Link{
				{
					SpanContext: apitrace.SpanContext{
						TraceID:    apitrace.ID{0xC0, 0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7, 0xC8, 0xC9, 0xCA, 0xCB, 0xCC, 0xCD, 0xCE, 0xCF},
						SpanID:     apitrace.SpanID{0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7},
						TraceFlags: 0,
					},
					Attributes: []label.KeyValue{
						label.String("LinkType", "Parent"),
					},
				},
				{
					SpanContext: apitrace.SpanContext{
						TraceID:    apitrace.ID{0xE0, 0xE1, 0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7, 0xE8, 0xE9, 0xEA, 0xEB, 0xEC, 0xED, 0xEE, 0xEF},
						SpanID:     apitrace.SpanID{0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7},
						TraceFlags: 0,
					},
					Attributes: []label.KeyValue{
						label.String("LinkType", "Child"),
					},
				},
			},
		*/
		// StatusCode:      codes.Internal,
		//StatusMessage:   "utterly unrecognized",
		HasRemoteParent: false,
		Attributes: []label.KeyValue{
			label.Int64("timeout_ns", 12e9),
		},
		//DroppedAttributeCount:    1,
		//DroppedMessageEventCount: 2,
		//DroppedLinkCount:         3,
		Resource: resource.New(label.String("rk1", "rv1"), label.Int64("rk2", 5)),
		InstrumentationLibrary: instrumentation.Library{
			Name:    "go.opentelemetry.io/test/otel",
			Version: "v0.0.1",
		},
	}

	got := transform.SpanData([]*export.SpanData{spanData})
	fmt.Printf("data size: %d\n", len(got))
	// 		_, err := e.traceExporter.Export(e.contextWithMetadata(ctx), &coltracepb.ExportTraceServiceRequest{
	//		ResourceSpans: protoSpans,
	//	})
	tr := coltracepb.ExportTraceServiceRequest{
		ResourceSpans: got,
	}
	b, err := json.Marshal(tr)
	if err != nil {
		fmt.Println("error:", err)
	}
	str = string(b)
	fmt.Println("str:", str)
	return str
}

func main() {
	url := "http://52.117.168.146:55681/v1/trace"
	method := "POST"

	var payload *strings.Reader

	if payloadType == "raw" {
		payload = strings.NewReader("{\n	  \"resource_spans\": [\n		{\n		  \"instrumentation_library_spans\": [\n			{\n			  \"spans\": [\n				{\n				  \"trace_id\": \"W47/95gDgQPSabYzgT/GDA==\",\n				  \"span_id\": \"7uGbfsPBsXM=\",\n				  \"name\": \"testSpan\",\n				  \"start_time_unix_nano\": 1544712660000000000,\n				  \"end_time_unix_nano\": 1544712661000000000,\n				  \"attributes\": [\n					{\n					  \"key\": \"attr1\",\n					  \"value\": { \"intValue\": 55 }\n					}\n				  ]\n				}\n			  ]\n			}\n		  ]\n		}\n	  ]\n	}")
	} else if payloadType == "protostr" {
		payload = strings.NewReader(makePayload())
	} else {
		fmt.Println("unknow type")
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		body, _ := ioutil.ReadAll(res.Body)

		fmt.Println(string(body))
		res.Body.Close()
	}
}
