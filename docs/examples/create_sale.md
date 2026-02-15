### Create Sale

This example demonstrates a semi-secure solution for creating a sale record, along with the associated
items and payment records. I say semi-secure because there isn't any verification of discount to 
sales, but products and discounts are loaded from the database and calculations are performed
server-side, then compared the expected total. 

This example demonstrates pulling records by id and inserting one record along with multiple
child records. Some mapper libraries could allow the 
developer to insert the parent and children in a single 
call, while some libaries on the other side require
setting the parent id on each child record individually.