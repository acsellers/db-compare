# Order Search

Pretty much any web application will need to support some sort of search functionality. 
At the moment, this does not have a full text search example, instead it is focused on 
matching on multiple columns across multiple tables with a single text like query.

This example is a simplified version of the sort of real-world query, but it has enough
complexity to demonstrate the strengths and weaknesses of the different mapper libraries.


### Example Implementation

```go
filters := []db.WhereFilter{}
if args.StartDate != "" {
    filters = append(filters, db.Orders.OrderDate.GTE(args.StartDate))
}
if args.EndDate != "" {
    filters = append(filters, db.Orders.OrderDate.LTE(args.EndDate))
}
if args.OrderType != "" {
    filters = append(filters, db.Orders.OrderType.EQ(args.OrderType))
}
if args.MinTotal != "" {
    filters = append(filters, db.Orders.Total.GTE(args.MinTotal))
}
if args.MaxTotal != "" {
    filters = append(filters, db.Orders.Total.LTE(args.MaxTotal))
}
if args.CustomerName != "" {
    filters = append(filters, db.Orders.Customers.Name.Like(args.CustomerName))
}

sales, err := db.Orders.Query(
    db.Orders.Joins.Customers,
    db.Orders.Where(filters...),
    db.Orders.OrderBy(db.Orders.OrderDate.Desc()),
)
```
