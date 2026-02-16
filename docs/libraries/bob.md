# Bob

Bob has a bit of everything as far as a database library goes. You give it
the database connection to load the schema, then it will generate a set of 
models for you. But you could also just use the query builder and the generic 
mapper code without using the model generation. It even has a query-based 
generated mapper, which is used in the Weekly Sales Report example.

On one hand, it's more Go code than sqlc, but you get some flexibility that 
isn't possible with sqlc. 

## Pros

* Table mapping is generated for the columns and does
not require reflection.
* Map tables and views, with relationships and type-specific query helpers.
* Has a fallback reflection-based struct mapper.
* Select specific columns, preload related tables, 

## Cons

* Need to be careful about generating the models when
feature branches have introduced new columns.