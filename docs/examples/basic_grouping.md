# Basic Grouping and Joining

There's three different examples in this section. One example just does
a group by, but the other two do joins as well a group by. The common package
refers to the common package in the store directory.

## Daily Revenue

This is the simplest example. It just groups the Orders table by date and 
order type and returns the total revenue for each.


### Example Implementation

```go
type DailyRevenue struct {
	OrderType    string    `json:"order_type" db:"order_type"`
	OrderDate    time.Time `json:"order_date" db:"order_date"`
	TotalRevenue float64   `json:"total_revenue" db:"total_revenue"`
}

results := []DailyRevenue{}

err := db.Orders.Query(
    db.GroupBy(db.Orders.OrderType, db.Orders.OrderDate),
    db.Select(
        db.Orders.OrderType,
        db.Orders.OrderDate,
        db.Orders.Total.Sum().As("total_revenue"),
    ),
).All(&results)
```

## Customer Sales

Now we add in a join to the Customers table and change to grouping by 
customer instead of order type. We also add in a count of the number of
orders.

### Example Implementation

```go
type CustomerTotal struct {
	ID          int64   `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	TotalSales  float64 `json:"total_sales" db:"total_sales"`
	TotalOrders int64   `json:"total_orders" db:"total_orders"`
}

results := []CustomerTotal{}

err := db.Customers.Query(
    db.Customers.Joins.Orders,
    db.GroupBy(db.Customers.ID, db.Customers.Name),
    db.Select(
        db.Customers.ID,
        db.Customers.Name,
        db.Customers.Orders.Count().As("order_count"),
        db.Customers.Orders.Total.Sum().As("total_sales"),
    ),
).All(&results)
```

## Daily Sold Items

This example goes back to not using joins, but it is using a View instead of
a Table. This is more a test of whether the mapper library can handle views.

### Example Implementation

```go
type ItemSummary struct {
	Name          string    `json:"name" db:"name"`
	Category      string    `json:"category" db:"category"`
	TotalQuantity int64     `json:"total_quantity" db:"total_quantity"`
	TotalSales    float64   `json:"total_sales" db:"total_sales"`
}

results := []ItemSummary{}

err := db.ItemSummary.Query(
    db.Where(db.ItemSummary.OrderDate.EQ(args.Date)),
    db.GroupBy(db.ItemSummary.Name, db.ItemSummary.Category),
    db.Select(
        db.ItemSummary.Name,
        db.ItemSummary.Category,
        db.ItemSummary.Quantity.Sum().As("total_quantity"),
        db.ItemSummary.Total.Sum().As("total_sales"),
    ),
).All(&results)
```
