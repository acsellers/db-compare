package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/acsellers/golang-db-compare/store/common"
	"github.com/acsellers/golang-db-compare/store/mysql/bob/models"
	"github.com/stephenafamo/bob/dialect/mysql"
)

func GetSale(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("id")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		fmt.Println(sid, err)
		http.Error(w, "Invalid sale ID", http.StatusBadRequest)
		return
	}
	order, err := models.Orders.Query(
		models.SelectWhere.Orders.ID.EQ(id),
		models.Preload.Order.Customer(
			mysql.PreloadOnly("name"),
		),
		models.SelectThenLoad.Order.OrderItems(
			models.Preload.OrderItem.Product(
				mysql.PreloadOnly("name", "category"),
			),
		),
		models.SelectThenLoad.Order.OrderPayments(),
	).One(r.Context(), db)
	if err != nil {
		fmt.Println("Query: ", err)
		http.Error(w, "Invalid sale ID", http.StatusBadRequest)
		return
	}
	sale := common.Sale{
		ID:             order.ID,
		OrderDate:      order.OrderDate,
		CustomerID:     order.CustomerID.Ptr(),
		DiscountID:     order.DiscountID.Ptr(),
		OrderType:      order.OrderType,
		Subtotal:       order.Subtotal.InexactFloat64(),
		DiscountAmount: order.DiscountAmount.InexactFloat64(),
		TaxAmount:      order.TaxAmount.InexactFloat64(),
		Total:          order.Total.InexactFloat64(),
		CreatedAt:      order.CreatedAt.GetOrZero(),
		UpdatedAt:      order.UpdatedAt.GetOrZero(),
	}
	if order.CustomerID.IsValue() && order.R.Customer != nil {
		sale.CustomerName = order.R.Customer.Name
	}
	for _, oi := range order.R.OrderItems {
		sale.Items = append(sale.Items, common.SaleItem{
			ID:              oi.ID,
			OrderID:         oi.OrderID,
			ProductID:       oi.ProductID,
			ProductName:     oi.R.Product.Name,
			ProductCategory: oi.R.Product.Category,
			DiscountID:      oi.DiscountID.Ptr(),
			Quantity:        int64(oi.Quantity),
			Price:           oi.Price.InexactFloat64(),
			DiscountAmount:  oi.DiscountAmount.InexactFloat64(),
			CreatedAt:       oi.CreatedAt.GetOrZero(),
			UpdatedAt:       oi.UpdatedAt.GetOrZero(),
		})
	}
	for _, op := range order.R.OrderPayments {
		info := map[string]any{}
		if op.PaymentInfo.IsValue() {
			json.Unmarshal(op.PaymentInfo.GetOrZero().Val, &info)
		}
		sale.Payments = append(sale.Payments, common.SalePayment{
			ID:          op.ID,
			OrderID:     op.OrderID,
			PaymentType: op.PaymentType,
			Amount:      op.Amount.InexactFloat64(),
			PaymentInfo: info,
			CreatedAt:   op.CreatedAt.GetOrZero(),
			UpdatedAt:   op.UpdatedAt.GetOrZero(),
		})
	}

	json.NewEncoder(w).Encode(sale)
}
