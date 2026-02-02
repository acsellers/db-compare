package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/acsellers/golang-db-compare/store/common"
)

func TestGetSales(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/sales/1")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	sale := common.Sale{}
	err = json.NewDecoder(resp.Body).Decode(&sale)
	if err != nil {
		t.Fatal(err)
	}
	if sale.ID != 1 {
		t.Errorf("Expected sale ID 1, got %d", sale.ID)
	}
	if sale.OrderDate.IsZero() {
		t.Error("Expected sale OrderDate to be set")
	}
	if sale.CustomerID == nil {
		t.Error("Expected sale CustomerID to be set")
	}
	if *sale.CustomerID != 1 {
		t.Errorf("Expected sale CustomerID 1, got %d", *sale.CustomerID)
	}
	if sale.CustomerName != "Alice Smith" {
		t.Errorf("Expected sale CustomerName 'Alice Smith', got %s", sale.CustomerName)
	}
	if len(sale.Items) != 1 {
		t.Errorf("Expected sale Items length 1, got %d", len(sale.Items))
	}
	if len(sale.Payments) != 2 {
		t.Errorf("Expected sale Payments length 2, got %d", len(sale.Payments))
	}
	if sale.Items[0].ProductName != "Smasnug Laptop 15 inch" {
		t.Errorf("Expected sale Items[0].ProductName 'Smasnug Laptop 15 inch', got %s", sale.Items[0].ProductName)
	}
	if sale.Items[0].ProductCategory != "pc" {
		t.Errorf("Expected sale Items[0].ProductCategory 'pc', got %s", sale.Items[0].ProductCategory)
	}
	if sale.Items[0].DiscountID != nil {
		t.Errorf("Expected sale Items[0].DiscountID to be nil")
	}
	if sale.Items[0].Quantity != 1 {
		t.Errorf("Expected sale Items[0].Quantity 1, got %d", sale.Items[0].Quantity)
	}
	if sale.Items[0].Price != 1200.00 {
		t.Errorf("Expected sale Items[0].Price 1200.00, got %f", sale.Items[0].Price)
	}
	if sale.Items[0].DiscountAmount != 0 {
		t.Errorf("Expected sale Items[0].DiscountAmount 0, got %f", sale.Items[0].DiscountAmount)
	}
	if sale.Items[0].CreatedAt.IsZero() {
		t.Error("Expected sale Items[0].CreatedAt to be set")
	}
	if sale.Items[0].UpdatedAt.IsZero() {
		t.Error("Expected sale Items[0].UpdatedAt to be set")
	}
	if sale.Items[0].ID != 1 {
		t.Errorf("Expected sale Items[0].ID 1, got %d", sale.Items[0].ID)
	}
	if sale.Items[0].OrderID != 1 {
		t.Errorf("Expected sale Items[0].OrderID 1, got %d", sale.Items[0].OrderID)
	}
	if sale.Items[0].ProductID != 1 {
		t.Errorf("Expected sale Items[0].ProductID 1, got %d", sale.Items[0].ProductID)
	}
	if sale.Items[0].ProductName != "Smasnug Laptop 15 inch" {
		t.Errorf("Expected sale Items[0].ProductName 'Smasnug Laptop 15 inch', got %s", sale.Items[0].ProductName)
	}
	if sale.Items[0].ProductCategory != "pc" {
		t.Errorf("Expected sale Items[0].ProductCategory 'pc', got %s", sale.Items[0].ProductCategory)
	}
	if sale.Items[0].DiscountID != nil {
		t.Errorf("Expected sale Items[0].DiscountID to be nil")
	}
	if sale.Items[0].Quantity != 1 {
		t.Errorf("Expected sale Items[0].Quantity 1, got %d", sale.Items[0].Quantity)
	}
	if sale.Items[0].Price != 1200.00 {
		t.Errorf("Expected sale Items[0].Price 1200.00, got %f", sale.Items[0].Price)
	}
	if sale.Items[0].DiscountAmount != 0 {
		t.Errorf("Expected sale Items[0].DiscountAmount 0, got %f", sale.Items[0].DiscountAmount)
	}
	if sale.Items[0].CreatedAt.IsZero() {
		t.Error("Expected sale Items[0].CreatedAt to be set")
	}
	if sale.Items[0].UpdatedAt.IsZero() {
		t.Error("Expected sale Items[0].UpdatedAt to be set")
	}
	if sale.Items[0].ID != 1 {
		t.Errorf("Expected sale Items[0].ID 1, got %d", sale.Items[0].ID)
	}
	if sale.Items[0].OrderID != 1 {
		t.Errorf("Expected sale Items[0].OrderID 1, got %d", sale.Items[0].OrderID)
	}
	if sale.Items[0].ProductID != 1 {
		t.Errorf("Expected sale Items[0].ProductID 1, got %d", sale.Items[0].ProductID)
	}
	if sale.Items[0].ProductName != "Smasnug Laptop 15 inch" {
		t.Errorf("Expected sale Items[0].ProductName 'Smasnug Laptop 15 inch', got %s", sale.Items[0].ProductName)
	}
	if sale.Items[0].ProductCategory != "pc" {
		t.Errorf("Expected sale Items[0].ProductCategory 'pc', got %s", sale.Items[0].ProductCategory)
	}
	if sale.Items[0].DiscountID != nil {
		t.Errorf("Expected sale Items[0].DiscountID to be nil")
	}
	if sale.Items[0].Quantity != 1 {
		t.Errorf("Expected sale Items[0].Quantity 1, got %d", sale.Items[0].Quantity)
	}
	if sale.Items[0].Price != 1200.00 {
		t.Errorf("Expected sale Items[0].Price 1200.00, got %f", sale.Items[0].Price)
	}
	if sale.Items[0].DiscountAmount != 0 {
		t.Errorf("Expected sale Items[0].DiscountAmount 0, got %f", sale.Items[0].DiscountAmount)
	}
	if sale.Items[0].CreatedAt.IsZero() {
		t.Error("Expected sale Items[0].CreatedAt to be set")
	}
	if sale.Items[0].UpdatedAt.IsZero() {
		t.Error("Expected sale Items[0].UpdatedAt to be set")
	}
}
