# Advanced Grouping and Joins

For now, there are just two examples of advanced grouping and joins.
Compared to some of the joins I write at work, these should actually
be in the basic section, but the current schema is pretty limited for 
showing off more complex joins.

Both examples use the ItemSummary view and join onto a ReportingOrders
table. The Typed Sales report also joins onto a dim_date table. Yes, it
is possible to use date functions to replace the dim_date table, but that
is out of scope for this example.

## General Sales Report

The General Sales report joins a view to a 
reporting category table, but only to the "general"
category labels.

### Example Implementation

```go
type GeneralSalesReport struct {
	Title       string  `json:"title" db:"title"`
	ReportOrder int64   `json:"report_order" db:"report_order"`
	ItemName    string  `json:"item_name" db:"item_name"`
	OrderCount  int64   `json:"order_count" db:"order_count"`
	Quantity    int64   `json:"quantity" db:"quantity"`
	TotalSales  float64 `json:"total_sales" db:"total_sales"`
}
results := []GeneralSalesReport{}

err := db.ItemSummaries.Query(
    db.Columns(
        models.ReportingOrders.Columns.Title,
        models.ReportingOrders.Columns.ReportOrder,
        models.ItemSummaries.Columns.Name.As("item_name"),
        models.ItemSummaries.Columns.OrderCount.Sum().As("order_count"),
        models.ItemSummaries.Columns.Quantity.Sum().As("quantity"),
        models.ItemSummaries.Columns.Total.Sum().As("total_sales"),
    ),
    db.ReportingOrder.JoinOn(
        db.ReportingOrder.Columns.OrderType.Eq("general"),
        db.ReportingOrder.Columns.Category.Eq(models.ItemSummaries.Columns.Category),
    ),
    db.ItemSummaries.OrderDate.GTE(startDate),
    db.ItemSummaries.OrderDate.LTE(endDate),
    db.GroupBy(
        db.ReportingOrders.Columns.Title,
        db.ReportingOrders.Columns.ReportOrder,
        db.ItemSummaries.Columns.Name,
    ),
    db.OrderBy(
        db.ReportingOrders.Columns.ReportOrder,
        db.ItemSummaries.Columns.Name,
    ),
).All(&results)
```

## Weekly Sales Report

The weekly (or typed) sales report steps up the complexity a bit by 
joining on a new dim_date table. For libraries that can do query
mapping and query building, they are allowed to use query mapping here, but 
not with the General Sales Report.

### Example Implementation

```go
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

results := []WeeklySaleReport{}

err := db.ItemSummaries.Query(
    db.Columns(
        models.DimDates.Columns.Year,
        models.DimDates.Columns.WeekOfYear,
        models.ReportingOrders.Columns.Title,
        models.ReportingOrders.Columns.ReportOrder,
        models.ItemSummaries.Columns.Name.As("item_name"),
        models.ItemSummaries.Columns.OrderCount.Sum().As("order_count"),
        models.ItemSummaries.Columns.Quantity.Sum().As("quantity"),
        models.ItemSummaries.Columns.Total.Sum().As("total_sales"),
    ),
    db.ReportingOrder.JoinOn(
        db.ReportingOrder.Columns.OrderType.Eq("general"),
        db.ReportingOrder.Columns.Category.Eq(models.ItemSummaries.Columns.Category),
    ),
    db.DimDate.JoinOn(
        db.DimDate.Columns.Date.Eq(models.ItemSummaries.Columns.OrderDate),
    ),
    db.ItemSummaries.OrderDate.GTE(startDate),
    db.ItemSummaries.OrderDate.LTE(endDate),
    db.GroupBy(
        db.DimDates.Columns.Year,
        db.DimDates.Columns.WeekOfYear,
        db.ReportingOrders.Columns.Title,
        db.ReportingOrders.Columns.ReportOrder,
        db.ItemSummaries.Columns.Name,
    ),
    db.OrderBy(
        db.DimDates.Columns.Year,
        db.DimDates.Columns.WeekOfYear,
        db.ReportingOrders.Columns.ReportOrder,
        db.ItemSummaries.Columns.Name,
    ),
).All(&results)
```