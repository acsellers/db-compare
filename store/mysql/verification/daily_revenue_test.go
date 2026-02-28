package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/acsellers/golang-db-compare/store/common"
)

func TestDailyRevenue(t *testing.T) {
	// [{"order_type":"members","order_date":"2026-01-01T00:00:00Z","total_revenue":3400.8},{"order_type":"non-members","order_date":"2026-01-01T00:00:00Z","total_revenue":605}]
	expected := map[string]common.DailyRevenue{
		"members": {
			TotalRevenue: 3400.8,
		},
		"non-members": {
			TotalRevenue: 605,
		},
	}
	resp, err := http.Get("http://localhost:8080/reports/daily-revenue?start_date=2026-01-01&end_date=2026-01-01")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	summs := []common.DailyRevenue{}
	if err := json.NewDecoder(resp.Body).Decode(&summs); err != nil {
		t.Fatal(err)
	}
	for _, s := range summs {
		expectation := expected[s.OrderType]
		if s.TotalRevenue != expectation.TotalRevenue {
			t.Errorf("expected %f, got %f", expectation.TotalRevenue, s.TotalRevenue)
		}
	}
}
