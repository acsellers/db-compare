package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/acsellers/golang-db-compare/store/common"
	"github.com/acsellers/golang-db-compare/store/mysql/sqlc/models"
)

// CreateSale creates a new sale from an Order struct. This involves
// retrieving the products and discounts, then calculating the totals and
// inserting the records into the database. If the expected total doesn't match
// the calculated total, it returns an error.
func CreateSale(w http.ResponseWriter, r *http.Request) {
	order := common.Order{}
	json.NewDecoder(r.Body).Decode(&order)

	// get discounts
	discounts := map[int64]models.Discount{}
	discountIds := []int64{}
	if order.DiscountID != nil {
		discountIds = append(discountIds, *order.DiscountID)
	}
	for _, oi := range order.Items {
		if oi.DiscountID != nil {
			discountIds = append(discountIds, *oi.DiscountID)
		}
	}
	if len(discounts) > 0 {
		discs, err := db.GetDiscount(r.Context(), discountIds)
		if err != nil {
			http.Error(w, "Invalid discount ID", http.StatusBadRequest)
			return
		}
		for _, d := range discs {
			discounts[d.ID] = d
		}
	}

	// get products
	products := map[int64]models.Product{}
	productIds := []int64{}
	for _, oi := range order.Items {
		productIds = append(productIds, oi.ProductID)
	}
	if len(productIds) == 0 {
		http.Error(w, "No products", http.StatusBadRequest)
		return
	}
	if len(productIds) > 0 {
		prods, err := db.GetProducts(r.Context(), productIds)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}
		for _, p := range prods {
			products[p.ID] = p
		}
	}

	// calculate totals
	subtotal := big.NewFloat(0)
	discountAmount := big.NewFloat(0)
	taxAmount := big.NewFloat(0)
	total := big.NewFloat(0)
	items := []models.CreateSaleItemsParams{}

	for _, oi := range order.Items {
		item := models.CreateSaleItemsParams{
			ProductID:  oi.ProductID,
			Quantity:   int32(oi.Quantity),
			Price:      products[oi.ProductID].Price,
			DiscountID: sql.NullInt64{Valid: oi.DiscountID != nil},
		}
		if oi.DiscountID != nil {
			item.DiscountID.Int64 = *oi.DiscountID
		}
		itemSubtotal := big.NewFloat(0.0)
		itemDiscountAmount := big.NewFloat(0.0)

		price := big.NewFloat(0.0)
		price.Parse(products[oi.ProductID].Price, 10)
		quantity := big.NewFloat(float64(oi.Quantity))
		itemSubtotal.Mul(price, quantity)
		subtotal.Add(subtotal, itemSubtotal)

		if oi.DiscountID != nil {
			discAmt := big.NewFloat(0.0)
			discAmt.Parse(discounts[*oi.DiscountID].Discount, 10)
			switch discounts[*oi.DiscountID].DiscountType {
			case "percentage":
				itemDiscountAmount.Mul(itemSubtotal, discAmt)
			case "fixed":
				itemDiscountAmount = discAmt
			}
		}
		discountAmount.Add(discountAmount, itemDiscountAmount)
		item.DiscountAmount = itemDiscountAmount.String()

		items = append(items, item)
	}

	if order.DiscountID != nil {
		discAmt := big.NewFloat(0.0)
		discAmt.Parse(discounts[*order.DiscountID].Discount, 10)
		switch discounts[*order.DiscountID].DiscountType {
		case "percentage":
			temp := big.NewFloat(0.0)
			discountAmount.Add(discountAmount, temp.Mul(subtotal, discAmt))
		case "fixed":
			discountAmount.Add(discountAmount, discAmt)
		}
	}

	total.Sub(subtotal, discountAmount)
	taxAmount.Mul(total, big.NewFloat(common.TaxRate))
	total.Add(total, taxAmount)

	if total.Text('f', 2) != fmt.Sprintf("%0.2f", order.ExpectedTotal) {
		http.Error(w, "Invalid total", http.StatusBadRequest)
		return
	}

	paymentSetters := []models.CreateSalePaymentsParams{}
	paymentAmount := big.NewFloat(0.0)
	for _, payment := range order.Payments {
		paymentAmount.Add(paymentAmount, big.NewFloat(payment.Amount))
		b, _ := json.Marshal(payment.PaymentInfo)
		paymentSetters = append(paymentSetters, models.CreateSalePaymentsParams{
			PaymentType: payment.PaymentType,
			Amount:      fmt.Sprintf("%0.2f", payment.Amount),
			PaymentInfo: b,
		})
	}
	if total.Cmp(paymentAmount) != 0 {
		http.Error(w, "Invalid payment amount", http.StatusBadRequest)
		return
	}

	orderType := "non-members"
	if order.CustomerID != nil {
		if exists, _ := db.CustomerExists(r.Context(), *order.CustomerID); exists == 0 {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			return
		}
		orderType = "members"
	}

	orderSetter := models.CreateSaleParams{
		CustomerID:     NullInt64FromPtr(order.CustomerID),
		DiscountID:     NullInt64FromPtr(order.DiscountID),
		OrderType:      orderType,
		Subtotal:       subtotal.Text('f', 2),
		DiscountAmount: discountAmount.Text('f', 2),
		TaxAmount:      taxAmount.Text('f', 2),
		Total:          total.Text('f', 2),
	}

	result, err := db.CreateSale(r.Context(), orderSetter)
	if err != nil {
		http.Error(w, "Invalid order", http.StatusBadRequest)
		return
	}
	orderId, _ := result.LastInsertId()
	for _, item := range items {
		item.OrderID = orderId
		err = db.CreateSaleItems(r.Context(), item)
		if err != nil {
			http.Error(w, "Invalid order items", http.StatusBadRequest)
			return
		}
	}
	for _, payment := range paymentSetters {
		payment.OrderID = orderId
		err = db.CreateSalePayments(r.Context(), payment)
		if err != nil {
			http.Error(w, "Invalid order payments", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(orderId)
}

func NullInt64FromPtr(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Valid: true, Int64: *i}
}
