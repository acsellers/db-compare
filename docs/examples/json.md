# JSON Support

JSON columns are a great way to store semi-structured data in a relational database. At the
same time, they're not well represented in the Go mapper library world. I would like to bring
more attention to this column type.

## JSON Saving/Loading

This evaluates whether the mapper library recognizes JSON columns and will add 
helpers to marshal/unmarshal the JSON into structs. 

### Example Record

```go
type OrderPayment struct {
	ID             int          `db:"id"`
	OrderID        int          `db:"order_id"`
	PaymentType    string       `db:"payment_type"`
	Amount         big.Float    `db:"amount"`
	PaymentInfo json.RawMessage `db:"payment_info"`
}
```

## JSON Querying

This evaluates whether the mapper library provides helpers to query JSON columns.

### Example Query

```go
now := time.Now()
expiringCards, err := db.OrderPayments.Query(
    db.OrderPayments.Where(
        db.OrderPayments.PaymentInfo.JSONPath("$.exp_month").Eq(int(now.Month())),
        db.OrderPayments.PaymentInfo.JSONPath("$.exp_year").Eq(now.Year()),
    )
)
```

## Auto Marshal/Unmarshal

Ideally, a mapper library would have some way for a developer to define the structure
of the column and be able to choose whether it automatically marshals/unmarshals the JSON into
structs. I don't know if any libraries currently do this, but I'll keep an eye out for it.

```go
type OrderPayment struct {
	ID             int             `db:"id"`
	OrderID        int             `db:"order_id"`
	PaymentType    string          `db:"payment_type"`
	Amount         big.Float       `db:"amount"`
	PaymentInfoRaw json.RawMessage `db:"payment_info"`
	PaymentInfo    PaymentInfo     `db:"payment_info"`
}

type PaymentInfo struct {
	Last4         string `json:"last4"`
	ExpMonth      int    `json:"exp_month"`
	ExpYear       int    `json:"exp_year"`
	CardType      string `json:"card_type"`
	TransactionID string `json:"transaction_id"`
}
```