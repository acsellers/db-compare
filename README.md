# Database Comparison

This is a set of common database operations that might be needed for a database mapping library.
Currently, it does not have a website, but that is the next step after I get more than one implementation
written. 

## Example Actions

### Create Sale

This takes in a sale and must retrieve products and discounts from the database to calculate the proper
total. Once the total is validated against the user-provided data, the sale is then created in the database
along with the sale items and discounts.

### Get Sale

Retrieve a sale record alongside some customer information and the order items and discounts.

### Bulk Load Customers

Given an uploaded file, load the customers into the database in the most efficient way possible. Note that
this is done through the mapper library, even if a more efficient way is possible by dropping down to the
database/sql library.

### Daily Revenue

This groups the Orders table by date and type, then calculates the sums of totals. It demonstrates a GROUP BY
clause in the mapper library.

### Daily Sold Items

This uses an Item Summaries view to calculate the total quantity sold and total sales for each item. It 
shows off the ability to use a view in the mapper library.

### General Sales

This uses the item summaries view and a oddly-joined table called reporting_order to calculate the total
quantity sold and total sales for each item. It then joins the item information to the reporting_order table
to format the data in a way that is easy to display. It demonstrates the use of a view and a join in the
mapper library.

### Typed Sales

This is very similar to the general sales report, but it joins on category and order type instead of just
order type like General Sales. Ideally, the join is able to have an on clause, but it can fallback to adding
where clauses to specify the join conditions.

## Mapper Libraries

### Bob

Bob is a schema-generated mapper libary. It passes all the requirements, though the General/Typed Sales has to
fallback to using where clauses to specify the join conditions. It also does not have a Bulk Insert helper beyond
a multi-insert helper. No Copy into or load data infile helpers are available.

### Ent

TODO

### sqlx

TODO

### GORM

TODO