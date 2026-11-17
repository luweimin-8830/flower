package handler

import (
	"bytes"
	"encoding/json"
	"errors"
)

// Omitted prices are preserved on updates; explicit null clears the price.
func parsePlantPrice(raw json.RawMessage) (*int64, error) {
	if len(raw) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return nil, nil
	}
	var cents int64
	if err := json.Unmarshal(raw, &cents); err != nil || cents < 0 || cents > 99999999 {
		return nil, errors.New("价格须为 0–999999.99 元，以整数分提交")
	}
	return &cents, nil
}
