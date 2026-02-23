# SQLC

sqlc is a query mapper tool that generates Go code from SQL queries
and schema files. It's well-supported and has a strong set of features.

It was also one of the easier libraries to setup for this project, since
I start by creating a schema and queries file. I just had to add the parameter
names to the queries to create the code.

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
