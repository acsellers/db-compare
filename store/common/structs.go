package common

import "time"

const TaxRate = 0.08

type Sale struct {
	ID             int64         `json:"id"`
	OrderDate      time.Time     `json:"order_date"`
	CustomerID     *int64        `json:"customer_id"`
	CustomerName   string        `json:"customer_name"`
	LocationID     int64         `json:"location_id"`
	LocationName   string        `json:"location_name"`
	DiscountID     *int64        `json:"discount_id"`
	OrderType      string        `json:"order_type"`
	Subtotal       float64       `json:"subtotal"`
	DiscountAmount float64       `json:"discount_amount"`
	TaxAmount      float64       `json:"tax_amount"`
	Total          float64       `json:"total"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	Items          []SaleItem    `json:"items"`
	Payments       []SalePayment `json:"payments"`
}

type SaleItem struct {
	ID              int64     `json:"id"`
	OrderID         int64     `json:"order_id"`
	ProductID       int64     `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductCategory string    `json:"product_category"`
	DiscountID      *int64    `json:"discount_id"`
	Quantity        int64     `json:"quantity"`
	Price           float64   `json:"price"`
	DiscountAmount  float64   `json:"discount_amount"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type SalePayment struct {
	ID          int64          `json:"id"`
	OrderID     int64          `json:"order_id"`
	PaymentType string         `json:"payment_type"`
	Amount      float64        `json:"amount"`
	PaymentInfo map[string]any `json:"payment_info"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type Order struct {
	CustomerID    *int64         `json:"customer_id"`
	DiscountID    *int64         `json:"discount_id"`
	LocationID    int64          `json:"location_id"`
	ExpectedTotal float64        `json:"expected_total"`
	Items         []OrderItem    `json:"items"`
	Payments      []OrderPayment `json:"payments"`
}
type OrderItem struct {
	ProductID  int64  `json:"product_id"`
	Quantity   int64  `json:"quantity"`
	DiscountID *int64 `json:"discount_id"`
}
type OrderPayment struct {
	PaymentType string           `json:"payment_type"`
	Amount      float64          `json:"amount"`
	PaymentInfo OrderPaymentInfo `json:"payment_info"`
}
type OrderPaymentInfo struct {
	CardNumber        string `json:"card_number"`
	CardType          string `json:"card_type"`
	CardExpiry        string `json:"card_expiry"`
	PurchaseOrder     string `json:"purchase_order"`
	TransactionID     string `json:"transaction_id"`
	TransactionResult string `json:"transaction_result"`
}
type ItemSummary struct {
	Name          string    `json:"name" db:"name"`
	Category      string    `json:"category" db:"category"`
	TotalQuantity int64     `json:"total_quantity" db:"total_quantity"`
	TotalSales    float64   `json:"total_sales" db:"total_sales"`
	OrderDate     time.Time `json:"order_date" db:"order_date"`
}
type DailyRevenue struct {
	OrderType    string    `json:"order_type" db:"order_type"`
	OrderDate    time.Time `json:"order_date" db:"order_date"`
	TotalRevenue float64   `json:"total_revenue" db:"total_revenue"`
}
type CustomerTotals struct {
	ID          int64   `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	TotalSales  float64 `json:"total_sales" db:"total_sales"`
	TotalOrders int64   `json:"total_orders" db:"total_orders"`
}

type GeneralSalesReport struct {
	Title       string  `json:"title" db:"title"`
	ReportOrder int64   `json:"report_order" db:"report_order"`
	ItemName    string  `json:"item_name" db:"item_name"`
	OrderCount  int64   `json:"order_count" db:"order_count"`
	Quantity    int64   `json:"quantity" db:"quantity"`
	TotalSales  float64 `json:"total_sales" db:"total_sales"`
}
type WeeklySaleReport struct {
	Year        int64   `json:"year" db:"year"`
	WeekOfYear  int64   `json:"week_of_year" db:"week_of_year"`
	Title       string  `json:"title" db:"title"`
	ReportOrder int64   `json:"report_order" db:"report_order"`
	ItemName    string  `json:"item_name" db:"item_name"`
	OrderCount  int64   `json:"order_count" db:"order_count"`
	Quantity    int64   `json:"quantity" db:"quantity"`
	TotalSales  float64 `json:"total_sales" db:"total_sales"`
}

type SaleSearch struct {
	CustomerName string  `json:"customer_name"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	OrderType    string  `json:"order_type"`
	MinTotal     float64 `json:"min_total"`
	MaxTotal     float64 `json:"max_total"`
	LocationID   int64   `json:"location_id"`
}

type SalesByCardType struct {
	OrderID    []int64 `json:"order_id" db:"order_id"`
	TotalSales float64 `json:"total_sales" db:"total_sales"`
	OrderCount int64   `json:"order_count" db:"order_count"`
}

type LocationPaymentInfo struct {
	LocationID          int64              `json:"location_id" db:"location_id"`
	TotalSales          float64            `json:"total_sales" db:"total_sales"`
	SalesByPaymentType  map[string]float64 `json:"sales_by_payment_type" db:"sales_by_payment_type"`
	SalesByCardType     map[string]float64 `json:"sales_by_card_type" db:"sales_by_card_type"`
	OrderCount          int64              `json:"order_count" db:"order_count"`
	OrdersByPaymentType map[string]int64   `json:"orders_by_payment_type" db:"orders_by_payment_type"`
	OrdersByCardType    map[string]int64   `json:"orders_by_card_type" db:"orders_by_card_type"`
}
