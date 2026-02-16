# Update Customer

This subject has two sub-examples, demonstrating two different
ways to update a bunch of customer records.



### Example Implementation

```go
for _, row := range inputRows {
    if row[0] == "" {
        continue
    }
    changes := []db.CustomerUpdateParams{}
    if row[1] != "" {
        changes = append(changes, db.Customers.Columns.Name.Set(row[1]))
    }
    if row[2] != "" {
        changes = append(changes, db.Customers.Columns.Email.Set(row[2]))
    }
    if row[3] != "" {
        changes = append(changes, db.Customers.Columns.Phone.Set(row[3]))
    }
    db.Customers.Update(
        db.Where(db.Customers.ExternalID.EQ(row[0])),
        changes...,
    ).Exec(r.Context())
}
```