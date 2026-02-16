# Get Sale

This is the first subject and arguable the simplest. That being said, not every library is going
to have the simplest solution. ORM style libraries that can preload the data will take the 
gold, while others will require multiple queries.

### Example Implementation

```go
sale, err := db.Sales.Query(
    db.Where(db.Sales.ID.EQ(args.ID)),
    db.Sales.Preload.OrderItems,
    db.Sales.Preload.OrderPayments,
).One(r.Context())
```
