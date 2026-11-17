package handler

import (
	"encoding/json"
	"testing"
)

func TestPlantPrice(t *testing.T) {
	for _, raw := range []string{"", "null"} {
		price, err := parsePlantPrice(json.RawMessage(raw))
		if err != nil || price != nil {
			t.Fatalf("empty %q: %v %v", raw, price, err)
		}
	}
	for _, raw := range []string{"0", "1", "12345", "99999999"} {
		price, err := parsePlantPrice(json.RawMessage(raw))
		if err != nil || price == nil {
			t.Fatalf("valid %q: %v", raw, err)
		}
	}
	for _, raw := range []string{"-1", "1.5", "100000000", `"123"`, "true", "{}"} {
		if _, err := parsePlantPrice(json.RawMessage(raw)); err == nil {
			t.Fatalf("accepted %s", raw)
		}
	}
	var omitted, cleared UpdatePlantRequest
	json.Unmarshal([]byte(`{"id":1}`), &omitted)
	json.Unmarshal([]byte(`{"id":1,"priceCents":null}`), &cleared)
	if len(omitted.PriceCents) != 0 || len(cleared.PriceCents) == 0 {
		t.Fatal("omitted and cleared must differ")
	}
}
