# SQLC

Let's be honest, this was the simplest library to setup
for this project. Every table and query was fully 
specified before implementation as part of planning.
The only changes I had to make was to add the names for
the query parameters.

## Pros

* If you're good at writing SQL, this is going to provide
a simple way to generate a good query library.
* Type specific query structs per query.
* Share structs within queries.

## Cons

* You're writing every single query, even the simple ones.
* You might need to write a lot of update queries or commit to updating all the columns.
* Some queries might need to be written in a fall-back
manner to handle optional parameters.
