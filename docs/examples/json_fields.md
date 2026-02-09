# JSON Fields

Ever since document databases became popular, developers have realized the power of 
stuffing json objects into columns. When you have a bunch of attributes that don't
really get queried, but need to be stored, you reach for a JSON column.

Ideally, the database mapper will automatically marshal and unmarshal the json 
struct in and out of the column. That's the ideal case, but it's not always present
in the mapper libraries.