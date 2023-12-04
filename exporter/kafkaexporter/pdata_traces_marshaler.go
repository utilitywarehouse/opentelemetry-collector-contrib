// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package kafkaexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter"

import (
	"github.com/IBM/sarama"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/traceutil"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchpersignal"
)

type pdataTracesMarshaler struct {
	marshaler ptrace.Marshaler
	encoding  string
}

func (p pdataTracesMarshaler) Marshal(td ptrace.Traces, topic string) ([]*sarama.ProducerMessage, error) {
	bts, err := p.marshaler.MarshalTraces(td)
	if err != nil {
		return nil, err
	}
	return []*sarama.ProducerMessage{
		{
			Topic: topic,
			Value: sarama.ByteEncoder(bts),
		},
	}, nil
}

func (p pdataTracesMarshaler) Encoding() string {
	return p.encoding
}

func (p pdataTracesMarshaler) KeyData() string {
	return "none"
}

type pdataTracesMarshalerByTraceId pdataTracesMarshaler

func (p pdataTracesMarshalerByTraceId) Marshal(td ptrace.Traces, topic string) ([]*sarama.ProducerMessage, error) {
	var messages []*sarama.ProducerMessage

	for _, tracesById := range batchpersignal.SplitTraces(td) {
		bts, err := p.marshaler.MarshalTraces(tracesById)
		if err != nil {
			return nil, err
		}
		var traceID = tracesById.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).TraceID()
		key := traceutil.TraceIDToHexOrEmptyString(traceID)
		messages = append(messages, &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(bts),
			Key:   sarama.ByteEncoder(key),
		})
	}
	return messages, nil
}

func (p pdataTracesMarshalerByTraceId) Encoding() string {
	return p.encoding
}

func (p pdataTracesMarshalerByTraceId) KeyData() string {
	return "traceID"
}

func newPdataTracesMarshaler(marshaler ptrace.Marshaler, encoding string) TracesMarshaler {
	return pdataTracesMarshaler{
		marshaler: marshaler,
		encoding:  encoding,
	}
}

func newPdataTracesMarshalerByTraceId(marshaler ptrace.Marshaler, encoding string) TracesMarshaler {
	return pdataTracesMarshalerByTraceId{
		marshaler: marshaler,
		encoding:  encoding,
	}
}
