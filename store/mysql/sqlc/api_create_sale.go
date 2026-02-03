package main

import (
	"database/sql"
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/acsellers/golang-db-compare/store/common"
	"github.com/acsellers/golang-db-compare/store/mysql/sqlc/models"
	"github.com/shopspring/decimal"
	"github.com/stephenafamo/bob/types"
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
			Price:      omit.From(products[oi.ProductID].Price),
			DiscountID: sql.NullInt64(oi.DiscountID),
		}
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
		item.DiscountAmount = omit.From(itemDiscountAmount)

		items = append(items, item)
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

	paymentSetters := []*models.OrderPaymentSetter{}
	paymentAmount := decimal.Zero
	for _, payment := range order.Payments {
		paymentAmount = paymentAmount.Add(decimal.NewFromFloat(payment.Amount))
		b, _ := json.Marshal(payment.PaymentInfo)
		paymentSetters = append(paymentSetters, &models.OrderPaymentSetter{
			PaymentType: omit.From(payment.PaymentType),
			Amount:      omit.From(decimal.NewFromFloat(payment.Amount)),
			PaymentInfo: omitnull.From(types.JSON[json.RawMessage]{Val: b}),
		})
	}
	if !total.Equal(paymentAmount) {
		http.Error(w, "Invalid payment amount", http.StatusBadRequest)
		return
	}

	orderType := "non-members"
	if order.CustomerID != nil {
		if exists, _ := models.CustomerExists(r.Context(), db, *order.CustomerID); !exists {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			return
		}
		orderType = "members"
	}

	orderSetter := &models.OrderSetter{
		CustomerID:     omitnull.FromPtr(order.CustomerID),
		DiscountID:     omitnull.FromPtr(order.DiscountID),
		OrderType:      omit.From(orderType),
		Subtotal:       omit.From(subtotal),
		DiscountAmount: omit.From(discountAmount),
		TaxAmount:      omit.From(taxAmount),
		Total:          omit.From(total),
	}

	orderRecord, err := models.Orders.Insert(orderSetter).One(r.Context(), db)
	if err != nil {
		http.Error(w, "Invalid order", http.StatusBadRequest)
		return
	}
	err = orderRecord.InsertOrderItems(r.Context(), db, items...)
	if err != nil {
		http.Error(w, "Invalid order items", http.StatusBadRequest)
		return
	}
	err = orderRecord.InsertOrderPayments(r.Context(), db, paymentSetters...)
	if err != nil {
		http.Error(w, "Invalid order payments", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(orderRecord.ID)
}
