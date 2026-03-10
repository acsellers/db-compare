# Library Checklist

## Docs

### Description.md

Write up a general description of the library and the features. Add
how simple it was to create the code for the examples. Add a pros and
cons section with a list of things that you like and dislike about the
library.

### Features.json

List database support, feature support and the activity/documentation/examples/migration
support of the library.

### Grades.json

Award grades for how well the library handled each set of example code.

### Info.json

Short info about the library, name, key, website, repo, a one line descriptiong,
which databases it supports, license, feature list, popularity (stars-ish).

### Samples.json

Use the sample editor to generate the samples.json file.

### Benchmarks.json

TODO: this reports how well the library performed on the benchmarks. Need to write
the benchmarks and ideally it will generate the benchmarks file.

## Examples

### 01 - Get Sale

#### A - Single Sale

Pull the sale record along with customer name (if present) and location name. Then add
the order items and order payments alongside the order. Use the common.Sale, common.SaleItem,
and common.SalePayment structs for the output.

GET /01/sales/1234

#### B - Location Sales for a Time Frame

Do the same as Single Sale, but return all sales for a location between two dates.

GET /01/sales?location=12&start_date=2025-01-01&end_date=2025-01-31

### 02 - Create Sale

Insert a sale, but validate that the sale items and discounts match up to an expected
number before inserting.

POST /02/sales {
    customer_id: 1234, (optional)
    discount_id: null, (optional)
    location_id: 123,
    expected_total: 20.00,
    items: [{
        product_id: 1234,
        quantity: 2,
        discount_id: null
    }],
    payments: [{
        payment_type: "cash",
        amount: 20.0,
        payment_info: {
            card_number: "",
            card_type: "",
            card_expiry: "",
            purchase_order: ""
        }
    }]
}

### 03 - Sale Search

Search sales based on a bunch of different fields, customer name (like), order date range,
order type, total greater than, total less than, or location id.

POST /03/sales/search {
    customer_name: "John",
    start_date: "2025-01-01",
    end_date: "2025-01-31",
    order_type: "members",
    min_total: 20.00,
    max_total: 100.00,
    location_id: 123
}

### 04 - Bulk Load Customers

What's the easiest way to insert a lot of customer records for this library. Test uses 10,000,
100,000, 1,000,000, and 10,000,000 records to see when the library fails (or if it never fails).
You cannot fall back to copy data or load data queries, it must use the provided code. 

### 05 - Customer Update

#### A - Individual Attribute Updates

In this case, do individual updates based on id. Ideally, the library can set the updated fields
dynamically, otherwise each field needs an individual query to update. Probably like 2,500
records to update.

POST /05/customers/update

Example CSV:

external_id,name,email,phone
1234,John Doe,john.doe@example.com,123-456-7890

#### B - Temporary Table Updates

Load data into a temporary table and run updates using that table. In this case, you can use
load data or copy into to fill the table. Think like 100,000 updates.

POST /05/customers/update/table

Example CSV:

external_id,name,email,phone
1234,John Doe,john.doe@example.com,123-456-7890

### 06 - JSON

#### A - Query

Get order ids and total sales (count and sum) by a customer for a specific brand of credit card,
optionally filtered by start and end date.

GET /06/customers/sales_by_card_type?card_type=Visa

#### B - Update

Update order payments in a webhook-like manner to set a field named TransactionID on them.

POST /06/webhook/update_payment
{
    order_id: 1234,
    transaction_id: "1234",
    transaction_result: "success"
}

#### C - Report

Look at payment info for a date range at a location, involves pulling back the order data and
parsing the payment info.

GET /06/locations/payments?id=123&start_date=2025-01-01&end_date=2025-01-31

### 07 - With Queries

### 08 - Basic Grouping

#### A - Daily Revenue

#### B - Customer Sales

#### C - Daily Sold Items

### 09 - Advanced Grouping

#### A - General Sales Report

#### B - Weekly Sales Report