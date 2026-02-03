package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/acsellers/golang-db-compare/store/common"
)

func toIntPtr(v sql.NullInt64) *int64 {
	if v.Valid {
		return &v.Int64
	}
	return nil
}
func toTimePtr(v sql.NullTime) *time.Time {
	if v.Valid {
		return &v.Time
	}
	return nil
}
func toStringPtr(v sql.NullString) *string {
	if v.Valid {
		return &v.String
	}
	return nil
}
func toStringOrZero(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}
func toBoolPtr(v sql.NullBool) *bool {
	if v.Valid {
		return &v.Bool
	}
	return nil
}
func parseStringFloat(v string) float64 {
	fv, _ := strconv.ParseFloat(v, 64)
	return fv
}
func parseNullStringFloat(v sql.NullString) float64 {
	if v.Valid {
		return parseStringFloat(v.String)
	}
	return 0
}
func GetSale(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("id")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		fmt.Println(sid, err)
		http.Error(w, "Invalid sale ID", http.StatusBadRequest)
		return
	}
	order, err := db.GetSale(r.Context(), id)
	if err != nil {
		fmt.Println("Query: ", err)
		http.Error(w, "Invalid sale ID", http.StatusBadRequest)
		return
	}
	sale := common.Sale{
		ID:             order.ID,
		OrderDate:      order.OrderDate,
		CustomerID:     toIntPtr(order.CustomerID),
		CustomerName:   toStringOrZero(order.Name),
		DiscountID:     toIntPtr(order.DiscountID),
		OrderType:      order.OrderType,
		Subtotal:       parseStringFloat(order.Subtotal),
		DiscountAmount: parseStringFloat(order.DiscountAmount),
		TaxAmount:      parseStringFloat(order.TaxAmount),
		Total:          parseStringFloat(order.Total),
		CreatedAt:      order.CreatedAt.Time,
		UpdatedAt:      order.UpdatedAt.Time,
	}
	items, err := db.GetSaleItems(r.Context(), id)
	if err != nil {
		fmt.Println("Query: ", err)
		http.Error(w, "Invalid sale ID", http.StatusBadRequest)
		return
	}
	for _, oi := range items {
		sale.Items = append(sale.Items, common.SaleItem{
			ID:              oi.ID,
			OrderID:         oi.OrderID,
			ProductID:       oi.ProductID,
			ProductName:     toStringOrZero(oi.Name),
			ProductCategory: toStringOrZero(oi.Category),
			DiscountID:      toIntPtr(oi.DiscountID),
			Quantity:        int64(oi.Quantity),
			Price:           parseStringFloat(oi.Price),
			DiscountAmount:  parseStringFloat(oi.DiscountAmount),
			CreatedAt:       oi.CreatedAt.Time,
			UpdatedAt:       oi.UpdatedAt.Time,
		})
	}
	payments, err := db.GetSalePayments(r.Context(), id)
	if err != nil {
		fmt.Println("Query: ", err)
		http.Error(w, "Invalid sale ID", http.StatusBadRequest)
		return
	}
	for _, op := range payments {
		info := map[string]any{}
		json.Unmarshal(op.PaymentInfo, &info)
		sale.Payments = append(sale.Payments, common.SalePayment{
			ID:          op.ID,
			OrderID:     op.OrderID,
			PaymentType: op.PaymentType,
			Amount:      parseStringFloat(op.Amount),
			PaymentInfo: info,
			CreatedAt:   op.CreatedAt.Time,
			UpdatedAt:   op.UpdatedAt.Time,
		})
	}

	json.NewEncoder(w).Encode(sale)
}
