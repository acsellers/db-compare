package main

import (
	"encoding/json"
	"net/http"

	"github.com/acsellers/golang-db-compare/store/common"
	"github.com/acsellers/golang-db-compare/store/mysql/bob/models"
	"github.com/shopspring/decimal"
)

// CreateSale creates a new sale from an Order struct. This involves
// retrieving the products and discounts, then calculating the totals and
// inserting the records into the database. If the expected total doesn't match
// the calculated total, it returns an error.
func CreateSale(w http.ResponseWriter, r *http.Request) {
	order := common.Order{}
	json.NewDecoder(r.Body).Decode(&order)

	// get discounts
	discounts := map[int64]*models.Discount{}
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
		// TODO: Add a query to get the discounts
	}

	// get products
	products := map[int64]*models.Product{}
	productIds := []int64{}
	for _, oi := range order.Items {
		productIds = append(productIds, oi.ProductID)
	}
	if len(productIds) == 0 {
		http.Error(w, "No products", http.StatusBadRequest)
		return
	}
	if len(productIds) > 0 {
		// TODO: Add a query to get the products
	}

	// calculate totals
	subtotal := decimal.Zero
	discountAmount := decimal.Zero
	taxAmount := decimal.Zero
	total := decimal.Zero
	items := []*models.OrderItemSetter{}

	for _, oi := range order.Items {
		// TODO: Setup insertion params for the order items (if needed)
		itemSubtotal := decimal.Zero
		itemDiscountAmount := decimal.Zero

		price := products[oi.ProductID].Price
		quantity := decimal.NewFromInt(oi.Quantity)
		itemSubtotal = price.Mul(quantity)
		subtotal = subtotal.Add(itemSubtotal)

		if oi.DiscountID != nil {
			discAmt := discounts[*oi.DiscountID].Discount
			switch discounts[*oi.DiscountID].DiscountType {
			case "percentage":
				itemDiscountAmount = price.Mul(quantity).Mul(discAmt)
			case "fixed":
				itemDiscountAmount = discAmt
			}
		}
		discountAmount = discountAmount.Add(itemDiscountAmount)
		// TODO: Set discount amount insertion param
	}

	if order.DiscountID != nil {
		discAmt := discounts[*order.DiscountID].Discount
		switch discounts[*order.DiscountID].DiscountType {
		case "percentage":
			discountAmount = discountAmount.Add(subtotal.Mul(discAmt))
		case "fixed":
			discountAmount = discountAmount.Add(discAmt)
		}
	}

	total = subtotal.Sub(discountAmount)
	taxAmount = total.Mul(decimal.NewFromFloat(common.TaxRate))
	total = total.Add(taxAmount)

	if !total.Equal(decimal.NewFromFloat(order.ExpectedTotal)) {
		http.Error(w, "Invalid total", http.StatusBadRequest)
		return
	}

	// TODO: Setup insertion params for the order payments (if needed)
	paymentAmount := decimal.Zero
	for _, payment := range order.Payments {
		paymentAmount = paymentAmount.Add(decimal.NewFromFloat(payment.Amount))
		// TODO: Set payment info insertion param
	}
	if !total.Equal(paymentAmount) {
		http.Error(w, "Invalid payment amount", http.StatusBadRequest)
		return
	}

	orderType := "non-members"
	if order.CustomerID != nil {
		// TODO: Check that the customer id exists
		// if exists, _ := models.CustomerExists(order.CustomerID); !exists {
		// 	http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		// 	return
		// }
		orderType = "members"
	}

	// TODO: Setup insertion params for the order (if needed)

	// TODO: Insert the order, order items, and order payments

	w.WriteHeader(http.StatusCreated)
	// TODO: Return the record ID
	//json.NewEncoder(w).Encode(record.ID)
}
