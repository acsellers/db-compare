# Bulk Customer Load

Most applications will need some way to load mass amounts of data into the database. 

This example is limited to whatever the mapper supports. The simplest solution is to just
run a bunch of insert statements, using a prepared statement. A better solution is to have
multiple rows per insert statement.

The gold standard is for the mapper to support the special bulk load statements, like LOAD DATA
INFILE for Mysql or COPY INTO for Postgres. 

This example has three sub-tests: a 10,000 row test, a 1,000,000 row test, and a 100,000,000 row test.


### Example Implementation

```go
rows := make([]db.Input)
for _, inputRow := range inputRows {
    rows = append(rows, db.Input{
        db.StringValue(inputRow.Name),
        db.NullStringValue(inputRow.Phone),
        db.NullStringValue(inputRow.Email),
    })
}

_, err := db.Customers.BulkInsert(
    db.Columns(db.Customers.Name, db.Customers.Email, db.Customers.Phone),
    rows,
)
```