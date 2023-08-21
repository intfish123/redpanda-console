// Copyright 2023 Redpanda Data, Inc.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.md
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0

package serde

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

var _ Serde = (*JsonSerde)(nil)

type JsonSerde struct{}

func (JsonSerde) Name() PayloadEncoding {
	return PayloadEncodingJSON
}

func (JsonSerde) DeserializePayload(record *kgo.Record, payloadType PayloadType) (RecordPayload, error) {
	payload := payloadFromRecord(record, payloadType)
	trimmed := bytes.TrimLeft(payload, " \t\r\n")

	return jsonDeserializePayload(trimmed)
}

func (JsonSerde) SerializeObject(obj any, payloadType PayloadType, opts ...SerdeOpt) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func jsonDeserializePayload(payload []byte) (RecordPayload, error) {
	trimmed := bytes.TrimLeft(payload, " \t\r\n")

	if len(trimmed) == 0 {
		return RecordPayload{}, fmt.Errorf("after trimming whitespaces there were no characters left")
	}

	startsWithJSON := trimmed[0] == '[' || trimmed[0] == '{'
	if !startsWithJSON {
		return RecordPayload{}, fmt.Errorf("first byte indicates this it not valid JSON, expected brackets")
	}

	var obj any
	err := json.Unmarshal(payload, &obj)
	if err != nil {
		return RecordPayload{}, fmt.Errorf("failed to parse JSON payload: %w", err)
	}

	return RecordPayload{
		NormalizedPayload:   payload,
		DeserializedPayload: obj,
		Encoding:            PayloadEncodingJSON,
	}, nil
}
