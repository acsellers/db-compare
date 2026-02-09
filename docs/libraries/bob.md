# Bob

Bob is a library that generates a custom query library
from the live schema of the database. 

## Pros

* Table mapping is generated for the columns and does
not require reflection.
* Map tables and views, with relationships and type-specific query helpers.
* Has a fallback reflection-based struct mapper.
* Select specific columns, preload related tables, 

## Cons

* Need to be careful about generating the models when
feature branches have introduced new columns.