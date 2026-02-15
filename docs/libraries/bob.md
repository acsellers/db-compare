# Bob

Bob is a set of libraries that can be as simple as a query builder, but is most
often used as a Code Generation ORM tool. At the same time, it can be used as a
query builder, a generic mapper, or even a query-based generated mapper.

## Pros

* Table mapping is generated for the columns and does
not require reflection.
* Map tables and views, with relationships and type-specific query helpers.
* Has a fallback reflection-based struct mapper.
* Select specific columns, preload related tables, 

## Cons

* Need to be careful about generating the models when
feature branches have introduced new columns.