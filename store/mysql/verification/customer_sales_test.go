package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/acsellers/golang-db-compare/store/common"
)

func TestCustomerSales(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/reports/customer-sales?start_date=2026-01-01&end_date=2026-01-03")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	totals := []common.CustomerTotals{}
	err = json.NewDecoder(resp.Body).Decode(&totals)
	if err != nil {
		t.Fatal(err)
	}
	expected := map[int]common.CustomerTotals{
		1: {
			ID:          1,
			Name:        "Alice Smith",
			TotalSales:  1308,
			TotalOrders: 1,
		},
		2: {
			ID:          2,
			Name:        "Bob Johnson",
			TotalSales:  2092.8,
			TotalOrders: 1,
		},
	}
	for _, total := range totals {
		expectation := expected[int(total.ID)]
		if expectation.ID != total.ID {
			t.Errorf("Expected %v, got %v", expectation, total)
		}
		if expectation.Name != total.Name {
			t.Errorf("Expected %v, got %v", expectation, total)
		}
		if expectation.TotalSales != total.TotalSales {
			t.Errorf("Expected %v, got %v", expectation, total)
		}
		if expectation.TotalOrders != total.TotalOrders {
			t.Errorf("Expected %v, got %v", expectation, total)
		}
	}
}
