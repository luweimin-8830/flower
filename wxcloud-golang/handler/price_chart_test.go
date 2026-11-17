package handler

import "testing"

func TestPriceSummaryBoundaries(t *testing.T) {
	rows := []priceCount{{Count: 3}}
	for _, p := range []int64{0, 4999, 5000, 9999, 10000, 19999, 20000, 49999, 50000, 99999, 100000, 99999999} {
		price := p
		rows = append(rows, priceCount{PriceCents: &price, Count: 2})
	}
	got := summarizePrices(rows)
	if got.PricedCount != 24 || got.UnpricedCount != 3 || got.TotalCents != 200739988 {
		t.Fatalf("unexpected summary: %+v", got)
	}
	for _, bucket := range got.Buckets {
		if bucket.Value != 4 {
			t.Fatalf("boundary assigned incorrectly: %+v", got.Buckets)
		}
	}
	empty := summarizePrices(nil)
	if empty.TotalCents != 0 || len(empty.Buckets) != 6 {
		t.Fatalf("empty: %+v", empty)
	}
}

func TestPriceDates(t *testing.T) {
	for _, dates := range [][2]string{{"2026-11-17", "2026-11-16"}, {"2026-02-30", "2026-03-01"}, {"", "2026-11-17"}} {
		if _, _, err := priceDateRange(dates[0], dates[1]); err == nil {
			t.Fatalf("accepted invalid dates: %v", dates)
		}
	}
	a, b, err := priceDateRange("2026-11-17", "2026-11-17")
	if err != nil || b.Sub(a).Hours() != 24 || b.Format("2006-01-02") != "2026-11-18" {
		t.Fatalf("inclusive day: %v %v %v", a, b, err)
	}
	_, offset := a.Zone()
	if offset != 8*3600 {
		t.Fatal("must use business timezone")
	}
}
