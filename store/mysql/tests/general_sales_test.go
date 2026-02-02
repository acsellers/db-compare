package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/acsellers/golang-db-compare/store/common"
)

func TestGeneralSales(t *testing.T) {
	const expectedSrc = `[{"title":"PC Items","report_order":1,"item_name":"Smasnug Laptop 15 inch","order_count":1,"quantity":1,"total_sales":1200},{"title":"Phone \u0026 Wearable Items","report_order":2,"item_name":"Smasnug Phone","order_count":1,"quantity":1,"total_sales":800},{"title":"Phone \u0026 Wearable Items","report_order":2,"item_name":"Smasnug Tablet","order_count":2,"quantity":3,"total_sales":1800},{"title":"Tax Items","report_order":4,"item_name":"Sales Tax","order_count":3,"quantity":3,"total_sales":335.8},{"title":"Discounts","report_order":5,"item_name":"$50 off","order_count":1,"quantity":1,"total_sales":-50},{"title":"Discounts","report_order":5,"item_name":"10% off","order_count":1,"quantity":1,"total_sales":-80},{"title":"Payments","report_order":6,"item_name":"Cash","order_count":1,"quantity":1,"total_sales":-608},{"title":"Payments","report_order":6,"item_name":"Credit Card","order_count":2,"quantity":2,"total_sales":-2697.8},{"title":"Payments","report_order":6,"item_name":"Gift Card","order_count":1,"quantity":1,"total_sales":-700}]`
	lines := []common.SaleReportLine{}
	json.Unmarshal([]byte(expectedSrc), &lines)
	lookup := map[string]common.SaleReportLine{}
	for _, line := range lines {
		lookup[line.Title+"|"+line.ItemName] = line
	}

	resp, err := http.Get("http://localhost:8080/reports/general-sales?start_date=2026-01-01&end_date=2026-01-31")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	respLines := []common.SaleReportLine{}
	json.NewDecoder(resp.Body).Decode(&respLines)

	checked := 0
	for _, line := range respLines {
		if _, ok := lookup[line.Title+"|"+line.ItemName]; !ok {
			t.Errorf("Missing line: %s", line.Title+"|"+line.ItemName)
		}
		expectedLine := lookup[line.Title+"|"+line.ItemName]
		if line.OrderCount != expectedLine.OrderCount {
			t.Errorf("Order count mismatch for %s: expected %d, got %d", line.Title+"|"+line.ItemName, expectedLine.OrderCount, line.OrderCount)
		}
		if line.Quantity != expectedLine.Quantity {
			t.Errorf("Quantity mismatch for %s: expected %d, got %d", line.Title+"|"+line.ItemName, expectedLine.Quantity, line.Quantity)
		}
		if line.TotalSales != expectedLine.TotalSales {
			t.Errorf("Total sales mismatch for %s: expected %f, got %f", line.Title+"|"+line.ItemName, expectedLine.TotalSales, line.TotalSales)
		}
		if line.TotalSales != expectedLine.TotalSales {
			t.Errorf("Total sales mismatch for %s: expected %f, got %f", line.Title+"|"+line.ItemName, expectedLine.TotalSales, line.TotalSales)
		}
		checked++
	}
	if checked != len(lines) {
		t.Errorf("Checked %d lines, expected %d", len(respLines), len(lines))
	}
}
