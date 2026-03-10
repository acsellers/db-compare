package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	OrderCount = 30000
	StartDate  = time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	EndDate    = time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
)

type LocationData struct {
	ID      int
	Name    string
	Address string
	City    string
	State   string
	Zip     string
	TaxRate float64
}

type ProductData struct {
	ID       int
	Name     string
	Category string
	Price    float64
}
type DiscountData struct {
	ID       int
	Name     string
	Category string
	Type     string
	Value    float64
}
type CustomerData struct {
	ID             int
	Name           string
	Email          string
	Phone          string
	ExternalID     string
	MarketingOptIn bool
	JoinLocationID int
}

var (
	// location to tax rate
	locationLookup = map[int]float64{}
	// product to price
	productLookup         = map[int]float64{}
	productCategoryLookup = map[int]string{}
	// customer to join location
	customerLookup = map[int]int{}
)

var src = rand.New(rand.NewSource(42))

func main() {
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	if os.Getenv("MODE") == "benchmark" {
		OrderCount = 30_000_000
		StartDate = time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)

	}
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM order_items")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM order_payments")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM orders")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM customers")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM locations")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM products")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM discounts")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM dim_date")
	if err != nil {
		log.Fatal(err)
	}

	LoadLocations(db)
	LoadProducts(db)
	LoadDiscounts(db)
	LoadCustomers(db)
	LoadDimDate(db)
	GenerateOrders(db, OrderCount)
}

func LoadLocations(db *sql.DB) {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM locations").Scan(&count)
	if count > 0 {
		fmt.Println("Locations already loaded")
		return
	}
	locations := []LocationData{}
	f, err := os.Open("locations.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.Read()
	id := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		taxRate, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			log.Fatal(err)
		}
		locationLookup[id] = taxRate
		locations = append(locations, LocationData{
			ID:      id,
			Name:    record[0],
			Address: record[1],
			City:    record[2],
			State:   record[3],
			Zip:     record[4],
			TaxRate: taxRate,
		})
		id++
	}
	query := "INSERT INTO locations (id, name, address, city, state, zip, tax_rate) VALUES "
	args := []any{}
	for _, location := range locations {
		query += "(?, ?, ?, ?, ?, ?, ?),"
		args = append(args, location.ID, location.Name, location.Address, location.City, location.State, location.Zip, location.TaxRate)
	}
	query = query[:len(query)-1]
	_, err = db.Exec(query, args...)
	if err != nil {
		fmt.Println(query)
		fmt.Println(args[0:25])
		log.Fatal(err)
	}
	fmt.Println("Loaded locations")
}

func LoadProducts(db *sql.DB) {
	products := []ProductData{}
	f, err := os.Open("products.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.Read()
	id := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		price, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		productLookup[id] = price
		productCategoryLookup[id] = record[1]
		products = append(products, ProductData{
			ID:       id,
			Name:     record[0],
			Category: record[1],
			Price:    price,
		})
		id++
	}
	query := "INSERT INTO products (id, name, category, price) VALUES "
	args := []any{}
	for _, product := range products {
		query += "(?, ?, ?, ?),"
		args = append(args, product.ID, product.Name, product.Category, product.Price)
	}
	query = query[:len(query)-1]
	_, err = db.Exec(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Loaded products")
}

func LoadDiscounts(db *sql.DB) {
	discounts := []DiscountData{}
	f, err := os.Open("discounts.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.Read()
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		value, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Fatal(err)
		}
		id, _ := strconv.Atoi(record[0])
		discounts = append(discounts, DiscountData{
			ID:       id,
			Name:     record[1],
			Category: record[2],
			Type:     record[3],
			Value:    value,
		})
	}
	query := "INSERT INTO discounts (id, name, category, discount_type, discount) VALUES "
	args := []any{}
	for _, discount := range discounts {
		query += "(?, ?, ?, ?, ?),"
		args = append(args, discount.ID, discount.Name, discount.Category, discount.Type, discount.Value)
	}
	query = query[:len(query)-1]
	_, err = db.Exec(query, args...)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadCustomers(db *sql.DB) {
	customers := []CustomerData{}
	f, err := os.Open("customers.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.Read()
	id := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		joinLocationID, _ := strconv.Atoi(record[4])
		customers = append(customers, CustomerData{
			ID:             id,
			Name:           record[0],
			Email:          record[2],
			Phone:          record[1],
			ExternalID:     fmt.Sprintf("C%06d", id),
			MarketingOptIn: record[3] == "1",
			JoinLocationID: joinLocationID,
		})
		customerLookup[id] = joinLocationID
		id++
		if id%1000 == 0 {
			query := "INSERT INTO customers (id, name, email, phone, external_id, marketing_opt_in, join_location_id) VALUES "
			args := []any{}
			for _, customer := range customers {
				query += "(?, ?, ?, ?, ?, ?, ?),"
				args = append(args, customer.ID, customer.Name, customer.Email, customer.Phone, customer.ExternalID, customer.MarketingOptIn, customer.JoinLocationID)
			}
			query = query[:len(query)-1]
			_, err = db.Exec(query, args...)
			if err != nil {
				log.Fatal(err)
			}
			customers = []CustomerData{}
		}
	}
	if len(customers) > 0 {
		query := "INSERT INTO customers (id, name, email, phone, external_id, marketing_opt_in, join_location_id) VALUES "
		args := []any{}
		for _, customer := range customers {
			query += "(?, ?, ?, ?, ?, ?, ?),"
			args = append(args, customer.ID, customer.Name, customer.Email, customer.Phone, customer.ExternalID, customer.MarketingOptIn, customer.JoinLocationID)
		}
		query = query[:len(query)-1]
		_, err = db.Exec(query, args...)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func LoadDimDate(db *sql.DB) {
	query := "INSERT INTO dim_date (date, month, year, quarter, day_of_week, day_of_month, day_of_year, week_of_year, week_of_month) VALUES (?, ?,?,?, ?,?,?, ?,?)"
	current := StartDate
	ending := time.Now().AddDate(1, 0, 0)
	yearWeek := 52
	monthWeek := 5
	for current.Before(ending) {
		if current.Weekday() == time.Monday {
			if current.Month() == time.January && current.Day() == 1 && yearWeek >= 10 {
				yearWeek = 1
			} else {
				yearWeek++
			}
			if current.Day() <= 7 {
				monthWeek = 1
			} else {
				monthWeek++
			}
		}
		quarter := 1
		if current.Month() >= time.April && current.Month() < time.July {
			quarter = 2
		} else if current.Month() >= time.July && current.Month() < time.October {
			quarter = 3
		} else if current.Month() >= time.October && current.Month() < time.January {
			quarter = 4
		}

		args := []any{
			current,
			current.Month(),
			current.Year(),
			quarter,
			current.Weekday(),
			current.Day(),
			current.YearDay(),
			yearWeek,
			monthWeek,
		}
		_, err := db.Exec(query, args...)
		if err != nil {
			log.Fatal(err)
		}
		current = current.AddDate(0, 0, 1)
	}
}
func GenerateOrders(db *sql.DB, count int) {
	current := StartDate

	perDay := 1 + count/int(EndDate.Sub(StartDate).Hours()/24)

	for current.Before(EndDate) {
		generateOrdersForDay(db, current, perDay)
		current = current.AddDate(0, 0, 1)
	}
	db.Exec("ALTER TABLE orders AUTO_INCREMENT = " + strconv.Itoa(orderId+1))
}
func generateOrdersForDay(db *sql.DB, date time.Time, count int) {
	for i := 0; i < count; i++ {
		ot := src.Intn(100)
		switch {
		case ot < 50:
			// few items
			generateNormalOrder(db, date)
		case ot < 75:
			// single item
			generateSingleItemOrder(db, date)
		case ot < 85:
			// bulk purchase
			generateBulkOrder(db, date)
		case ot < 95:
			// discounted purchase
			generateDiscountedOrder(db, date)
		case ot < 98:
			// special discount
			generateSpecialDiscountOrder(db, date)
		default:
			// big order (0.2%)
			if src.Intn(10) < 2 {
				generateBigOrder(db, date)
			} else {
				generateDiscountedOrder(db, date)
			}
		}
	}
}

type orderData struct {
	id             int
	orderDate      time.Time
	customerID     *int
	discountID     *int
	orderType      string
	subtotal       float64
	discountAmount float64
	taxAmount      float64
	total          float64
	locationID     int
	orderItems     []orderItemData
	orderPayments  []orderPaymentData
}
type orderItemData struct {
	productID      int
	quantity       int
	price          float64
	discountID     *int
	discountAmount float64
	itemTotal      float64
}
type orderPaymentData struct {
	paymentType string
	amount      float64
	paymentInfo json.RawMessage
}

var orderId = 1

func genOrderData(date time.Time) *orderData {
	od := &orderData{
		id:        orderId,
		orderDate: date.Add(time.Second * time.Duration(src.Intn(24*60*60))),
	}
	if src.Intn(10) < 5 {
		ci := src.Intn(len(customerLookup)) + 1
		od.customerID = &ci
		if src.Intn(10) < 8 {
			od.locationID = customerLookup[ci-1]
		} else {
			od.locationID = src.Intn(len(locationLookup)) + 1
		}
		if od.locationID == 0 {
			od.locationID = src.Intn(len(locationLookup)) + 1
		}
	} else {
		od.locationID = src.Intn(len(locationLookup)) + 1
	}
	orderId++
	return od
}
func (od *orderData) addItem(id int, qty int) {
	price := productLookup[id]
	category := productCategoryLookup[id]
	if od.orderType == "" {
		od.orderType = category
	} else {
		if od.orderType != category {
			od.orderType = "mixed"
		}
	}
	oi := orderItemData{
		productID: id,
		quantity:  qty,
		price:     price,
		itemTotal: price * float64(qty),
	}
	od.orderItems = append(od.orderItems, oi)
	od.subtotal += oi.itemTotal
}
func randomCardNumber() string {
	return fmt.Sprintf("%016d", src.Int63())
}
func randomCardType() string {
	return []string{"Visa", "Mastercard", "American Express", "Discover"}[src.Intn(4)]
}
func randomCardExpiry() string {
	return fmt.Sprintf("%02d/%02d", src.Intn(12)+1, src.Intn(24)+25)
}
func (od *orderData) addPayment(paymentType string, amount float64) {
	switch paymentType {
	case "credit_card":
		od.orderPayments = append(od.orderPayments, orderPaymentData{
			paymentType: paymentType,
			amount:      amount,
			paymentInfo: []byte(`{"card_number": "` + randomCardNumber() + `", "card_type": "` + randomCardType() + `", "card_expiry": "` + randomCardExpiry() + `"}`),
		})
	case "purchase_order":
		od.orderPayments = append(od.orderPayments, orderPaymentData{
			paymentType: paymentType,
			amount:      amount,
			paymentInfo: []byte(`{"purchase_order": "` + fmt.Sprintf("%06d", src.Intn(1000000)) + `"}`),
		})
	case "gift_card":
		od.orderPayments = append(od.orderPayments, orderPaymentData{
			paymentType: paymentType,
			amount:      amount,
			paymentInfo: []byte(`{"gift_card": "` + fmt.Sprintf("%06d", src.Intn(1000000)) + `"}`),
		})
	case "cash":
		od.orderPayments = append(od.orderPayments, orderPaymentData{
			paymentType: paymentType,
			amount:      amount,
			paymentInfo: []byte("{}"),
		})
	}

}
func (od *orderData) calculateTotal() {
	od.taxAmount = od.subtotal * locationLookup[od.locationID]
	od.total = od.subtotal + od.taxAmount
}

var (
	pendingOrderCount        int
	pendingOrderItemCount    int
	pendingOrderPaymentCount int
	// id, date, customer_id, discount_id, order_type, subtotal, discount_amount, tax_amount, total, location_id
	orderArgs []any
	// order_id, product_id, quantity, discount_id, discount_amount
	orderItemArgs []any
	// order_id, payment_type, amount, payment_info
	orderPaymentArgs []any
)

func (od *orderData) insertOrder(db *sql.DB) {
	if pendingOrderCount == 1000 {
		insertPendingOrders(db)
	}
	orderArgs = append(orderArgs, od.id, od.orderDate, od.customerID, od.discountID, od.orderType, od.subtotal, od.discountAmount, od.taxAmount, od.total, od.locationID)
	for _, oi := range od.orderItems {
		pendingOrderItemCount++
		orderItemArgs = append(orderItemArgs, od.id, oi.productID, oi.quantity, oi.discountID, oi.discountAmount, oi.price, oi.itemTotal)
	}
	for _, op := range od.orderPayments {
		pendingOrderPaymentCount++
		orderPaymentArgs = append(orderPaymentArgs, od.id, op.paymentType, op.amount, op.paymentInfo)
	}
	pendingOrderCount++
}
func insertPendingOrders(db *sql.DB) {
	orderQuery := "INSERT INTO orders (id, order_date, customer_id, discount_id, order_type, subtotal, discount_amount, tax_amount, total, location_id) VALUES "
	for i := 0; i < pendingOrderCount; i++ {
		orderQuery += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"
	}
	orderQuery = orderQuery[:len(orderQuery)-1]
	_, err := db.Exec(orderQuery, orderArgs...)
	if err != nil {
		fmt.Println(orderQuery)
		fmt.Println(orderArgs)
		log.Fatal(err)
	}

	orderItemQuery := "INSERT INTO order_items (order_id, product_id, quantity, discount_id, discount_amount, price, item_total) VALUES "
	for i := 0; i < pendingOrderItemCount; i++ {
		orderItemQuery += "(?, ?, ?, ?, ?, ?, ?),"
	}
	orderItemQuery = orderItemQuery[:len(orderItemQuery)-1]
	_, err = db.Exec(orderItemQuery, orderItemArgs...)
	if err != nil {
		fmt.Println(orderItemQuery)
		fmt.Println(orderItemArgs)
		log.Fatal(err)
	}

	orderPaymentQuery := "INSERT INTO order_payments (order_id, payment_type, amount, payment_info) VALUES "
	for i := 0; i < pendingOrderPaymentCount; i++ {
		orderPaymentQuery += "(?, ?, ?, ?),"
	}
	orderPaymentQuery = orderPaymentQuery[:len(orderPaymentQuery)-1]
	_, err = db.Exec(orderPaymentQuery, orderPaymentArgs...)
	if err != nil {
		fmt.Println(orderPaymentQuery)
		fmt.Println(orderPaymentArgs)
		log.Fatal(err)
	}

	orderArgs = []any{}
	orderItemArgs = []any{}
	orderPaymentArgs = []any{}
	pendingOrderCount = 0
	pendingOrderItemCount = 0
	pendingOrderPaymentCount = 0
}
func generateNormalOrder(db *sql.DB, date time.Time) {
	od := genOrderData(date)
	for i := 0; i < src.Intn(3)+2; i++ {
		od.addItem(src.Intn(len(productLookup))+1, src.Intn(5)+1)
	}
	od.calculateTotal()
	paymentAmount := od.total
	if src.Intn(10) == 0 && od.total > 100 {
		amount := od.total * float64(src.Intn(100)) / 100.0
		od.addPayment("gift_card", amount)
		paymentAmount -= amount
	}
	if src.Intn(5) == 0 {
		od.addPayment("cash", paymentAmount)
	} else {
		od.addPayment("credit_card", paymentAmount)
	}
	od.insertOrder(db)
}
func generateSingleItemOrder(db *sql.DB, date time.Time) {
	od := genOrderData(date)
	od.addItem(src.Intn(len(productLookup))+1, 1)
	od.calculateTotal()
	if src.Intn(5) == 0 {
		od.addPayment("cash", od.total)
	} else {
		od.addPayment("credit_card", od.total)
	}
	od.insertOrder(db)
}
func generateBulkOrder(db *sql.DB, date time.Time) {
	od := genOrderData(date)
	for i := 0; i < src.Intn(12)+2; i++ {
		od.addItem(src.Intn(len(productLookup))+1, src.Intn(500)+1)
	}
	od.calculateTotal()
	od.addPayment("credit_card", od.total)
	od.insertOrder(db)
}
func getFixedDiscount(special bool, price float64, count int) (int, float64) {
	if special {
		return 21, 100.00
	}
	id := src.Intn(3) + 1
	switch id {
	case 1:
		return 1, 10.00
	case 2:
		return 2, 20.00
	case 3:
		return 3, 5.00
	}
	return 0, 0
}
func getPercentDiscount(special bool, price float64, count int) (int, float64) {
	if special {
		return 22, 0.50
	}
	id := src.Intn(3)
	amount := price * float64(count)
	switch id {
	case 0:
		return 11, amount * 0.05
	case 1:
		return 12, amount * 0.10
	case 2:
		return 13, amount * 0.15
	}
	return 0, 0
}
func generateDiscountedOrder(db *sql.DB, date time.Time) {
	od := genOrderData(date)
	itemDiscount := src.Intn(3) == 0
	for i := 0; i < src.Intn(3)+2; i++ {
		od.addItem(src.Intn(len(productLookup))+1, src.Intn(5)+1)
		if itemDiscount && od.orderItems[len(od.orderItems)-1].itemTotal > 25 {
			if src.Intn(3) == 0 {
				did, amt := getFixedDiscount(false, od.orderItems[len(od.orderItems)-1].itemTotal, 0)
				od.orderItems[len(od.orderItems)-1].discountID = &did
				od.orderItems[len(od.orderItems)-1].discountAmount = amt
				od.orderItems[len(od.orderItems)-1].itemTotal -= amt
				od.subtotal -= amt
			} else {
				did, amt := getPercentDiscount(false, od.orderItems[len(od.orderItems)-1].itemTotal, 0)
				od.orderItems[len(od.orderItems)-1].discountID = &did
				od.orderItems[len(od.orderItems)-1].discountAmount = amt
				od.orderItems[len(od.orderItems)-1].itemTotal -= amt
				od.subtotal -= amt
			}
		}
	}
	if !itemDiscount {
		// percent discount
		if src.Intn(3) == 0 {
			did, amt := getPercentDiscount(false, od.subtotal, 1)
			od.discountID = &did
			od.discountAmount = amt
		} else {
			// fixed discount
			did, amt := getFixedDiscount(false, od.subtotal, 1)
			od.discountID = &did
			od.discountAmount = amt
		}
	}
	od.calculateTotal()
	paymentAmount := od.total
	if src.Intn(10) == 0 && od.total > 50 {
		amount := 10.0 * float64(src.Intn(5)+1)
		od.addPayment("gift_card", amount)
		paymentAmount -= amount
	}
	if src.Intn(5) == 0 {
		od.addPayment("cash", paymentAmount)
	} else {
		od.addPayment("credit_card", paymentAmount)
	}
	od.insertOrder(db)
}
func generateSpecialDiscountOrder(db *sql.DB, date time.Time) {
	od := genOrderData(date)
	itemDiscount := src.Intn(3) == 0
	id := src.Intn(len(productLookup)) + 1
	qty := src.Intn(25) + int(100/productLookup[id])
	od.addItem(id, qty)
	if itemDiscount && od.orderItems[len(od.orderItems)-1].price > 25 {
		if src.Intn(3) == 0 {
			did, amt := getFixedDiscount(true, od.orderItems[len(od.orderItems)-1].price, 0)
			od.orderItems[len(od.orderItems)-1].discountID = &did
			od.orderItems[len(od.orderItems)-1].discountAmount = amt
		} else {
			did, amt := getPercentDiscount(true, od.orderItems[len(od.orderItems)-1].price, 0)
			od.orderItems[len(od.orderItems)-1].discountID = &did
			od.orderItems[len(od.orderItems)-1].discountAmount = amt
		}
	}
	for i := 0; i < src.Intn(3)+2; i++ {
		od.addItem(src.Intn(len(productLookup))+1, src.Intn(15)+1)
	}

	if !itemDiscount {
		// percent discount
		if src.Intn(3) == 0 {
			did, amt := getPercentDiscount(true, od.subtotal, 1)
			od.discountID = &did
			od.discountAmount = amt
		} else {
			// fixed discount
			did, amt := getFixedDiscount(true, od.subtotal, 1)
			od.discountID = &did
			od.discountAmount = amt
		}
	}
	od.calculateTotal()
	paymentAmount := od.total
	if src.Intn(10) == 0 && od.total > 50 {
		amount := od.total * float64(src.Intn(100)) / 100.0
		od.addPayment("gift_card", amount)
		paymentAmount -= amount
	}
	if src.Intn(5) == 0 {
		od.addPayment("cash", paymentAmount)
	} else {
		od.addPayment("credit_card", paymentAmount)
	}
	od.insertOrder(db)
}
func generateBigOrder(db *sql.DB, date time.Time) {
	od := genOrderData(date)
	for i := 0; i < src.Intn(12)+2; i++ {
		od.addItem(src.Intn(len(productLookup))+1, src.Intn(5_000)+10_000)
	}
	od.calculateTotal()
	od.addPayment("purchase_order", od.total)
	od.insertOrder(db)
}
