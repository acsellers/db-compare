# Create Sale

This example demonstrates a semi-secure solution for creating a sale record, along with the associated
items and payment records. I say semi-secure because there isn't any verification of discount to 
sales, but products and discounts are loaded from the database and calculations are performed
server-side, then compared the expected total. 

This example demonstrates pulling records by id and inserting one record along with multiple
child records. Some mapper libraries could allow the 
developer to insert the parent and children in a single 
call, while some libaries on the other side require
setting the parent id on each child record individually.

### Example Implementation (Abridged)

```go
discountIds := []int64{}
productIds := []int64{}
customerID := sql.NullInt64{}

// Find the discounts and products used in the sale...

discounts, err := db.Discounts.Query(
    db.Where(db.Discounts.ID.In(args.DiscountIDs...)),
).All(r.Context())

products, err := db.Products.Query(
    db.Where(db.Products.ID.In(args.ProductIDs...)),
).All(r.Context())

// Validate the order and setup the items and payments...

order := models.OrderParams{
    CustomerID: customerID,
    // Order Attributes: CustomerID, OrderType...
    OrderItems: orderItems,
    OrderPayments: orderPayments,
}

result, err := db.Orders.Insert(order).Exec(r.Context())
```