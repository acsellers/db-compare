package common

import "time"

const TaxRate = 0.08

type Sale struct {
	ID             int64         `json:"id"`
	OrderDate      time.Time     `json:"order_date"`
	CustomerID     *int64        `json:"customer_id"`
	CustomerName   string        `json:"customer_name"`
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
	PaymentType string         `json:"payment_type"`
	Amount      float64        `json:"amount"`
	PaymentInfo map[string]any `json:"payment_info"`
}
type ItemSummary struct {
	Name          string    `json:"name"`
	Category      string    `json:"category"`
	TotalQuantity int64     `json:"total_quantity"`
	TotalSales    float64   `json:"total_sales"`
	OrderDate     time.Time `json:"order_date"`
}
type DailyRevenue struct {
	OrderType    string    `json:"order_type"`
	OrderDate    time.Time `json:"order_date"`
	TotalRevenue float64   `json:"total_revenue"`
}
type CustomerTotals struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	TotalSales  float64 `json:"total_sales"`
	TotalOrders int64   `json:"total_orders"`
}

type SaleReportLine struct {
	Title       string  `json:"title"`
	ReportOrder int64   `json:"report_order"`
	ItemName    string  `json:"item_name"`
	OrderCount  int64   `json:"order_count"`
	Quantity    int64   `json:"quantity"`
	TotalSales  float64 `json:"total_sales"`
}
type WeeklySaleReport struct {
	Year        int64   `json:"year"`
	WeekOfYear  int64   `json:"week_of_year"`
	Title       string  `json:"title"`
	ReportOrder int64   `json:"report_order"`
	ItemName    string  `json:"item_name"`
	OrderCount  int64   `json:"order_count"`
	Quantity    int64   `json:"quantity"`
	TotalSales  float64 `json:"total_sales"`
}
