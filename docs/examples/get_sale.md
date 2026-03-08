# Get Sale


## Single Sale

This is the first subject and arguable the simplest. That being said, not every library is going
to have the simplest solution. ORM style libraries that can preload the data will take the 
gold, while others that require multiple queries will take lower marks.

### Example Implementation

```go
sale, err := db.Sales.Query(
    db.Where(db.Sales.ID.EQ(args.ID)),
    db.Sales.Preload.OrderItems,
    db.Sales.Preload.OrderPayments,
    db.Sales.Preload.Location,
    db.Sales.Preload.Customer,
).One(r.Context())
```

## Customer Sales for a year

This is basically the same as the single sale, but it rewards libraries that can 
preload relationships a bit more. You're given a customer id and a year and you need
to return all of the sales for that customer for that year.

### Example Implementation

```go

```
