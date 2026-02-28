package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/acsellers/golang-db-compare/store/common"
)

func TestDailySoldItems(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/reports/daily-sold-items?date=2026-01-01")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	summaries := []common.ItemSummary{}
	err = json.NewDecoder(resp.Body).Decode(&summaries)
	if err != nil {
		t.Fatal(err)
	}
	// non-exhaustive
	expected := map[string]common.ItemSummary{
		"Smasnug Laptop 15 inch": {
			Name:          "Smasnug Laptop 15 inch",
			Category:      "pc",
			TotalQuantity: 1,
			TotalSales:    1200,
		},
		"Smasnug Tablet": {
			Name:          "Smasnug Tablet",
			Category:      "phone",
			TotalQuantity: 3,
			TotalSales:    1800,
		},
	}

	checked := 0
	for _, summary := range summaries {
		expectation, ok := expected[summary.Name]
		if !ok {
			continue
		}
		checked++
		if expectation.Name != summary.Name {
			t.Errorf("Expected %v, got %v", expectation, summary)
		}
		if expectation.Category != summary.Category {
			t.Errorf("Expected %v, got %v", expectation, summary)
		}
		if expectation.TotalQuantity != summary.TotalQuantity {
			t.Errorf("Expected %v, got %v", expectation, summary)
		}
		if expectation.TotalSales != summary.TotalSales {
			t.Errorf("Expected %v, got %v", expectation, summary)
		}
	}
	if checked != len(expected) {
		t.Errorf("Expected %v, got %v", len(expected), checked)
	}

}
