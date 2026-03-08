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

#### B - Location Sales for a Time Frame

Do the same as Single Sale, but return all sales for a location between two dates.

### 02 - Create Sale

Insert a sale, but validate that the sale items and discounts match up to an expected
number before inserting.

### 03 - Sale Search

Search sales based on a bunch of different fields, customer name (like), order date range,
order type, total greater than, total less than, or location id.

### 04 - Bulk Load Customers

What's the easiest way to insert a lot of customer records for this library. Test uses 10,000,
100,000, 1,000,000, and 10,000,000 records to see when the library fails (or if it never fails).
You cannot fall back to copy data or load data queries, it must use the provided code. 

### 05 - Customer Update

#### A - Individual Attribute Updates

In this case, do individual updates based on id. Ideally, the library can set the updated fields
dynamically, otherwise each field needs an individual query to update. Probably like 2,500
records to update.

#### B - Temporary Table Updates

Load data into a temporary table and run updates using that table. In this case, you can use
load data or copy into to fill the table. Think like 100,000 updates.

### 06 - JSON

#### A - Query

Query all sales by a customer to filter by type of credit card.

#### B - Update

Update order payments in a webhook-like manner to set a field named TransactionID on them.

#### C - Report

Look at payment info for a date range at a location, involves pulling back the order data and
parsing the payment info.

### 07 - With Queries

### 08 - Basic Grouping

#### A - Daily Revenue

#### B - Customer Sales

#### C - Daily Sold Items

### 09 - Advanced Grouping

#### A - General Sales Report

#### B - Weekly Sales Report