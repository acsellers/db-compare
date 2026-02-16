# Update Customer

While the customer record in this schema doesn't have 
many columns, most customer tables will have plenty of
different columns to update. This example demonstrates
how the mapping library handles updating specific 
columns.

Does the mapping library really only have single column
or all column upates? Or can the mapping library update
columns dynamically?

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