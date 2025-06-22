package gox

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

func GobMarshal[T any](v T) []byte {
	var buf bytes.Buffer
	_ = gob.NewEncoder(&buf).Encode(v)
	return buf.Bytes()
}

func GobUnmarshal[T any](b []byte) (T, error) {
	var v T
	return v, gob.NewDecoder(bytes.NewBuffer(b)).Decode(&v)
}

func Marshal[T any](v T) []byte {
	b, _ := json.Marshal(v)
	return b
}

func Unmarshal[T any](b []byte) (T, error) {
	var v T
	return v, json.Unmarshal(b, &v)
}
