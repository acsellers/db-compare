-- get discount
SELECT * FROM discounts WHERE id IN (?);

-- get products
SELECT * FROM products WHERE id IN (?);

-- get sale
select orders.*, customers.name
from orders 
left join customers on orders.customer_id = customers.id
where orders.id = ?;
select order_items.*, 
  products.name, products.category 
from order_items 
left join products on order_items.product_id = products.id
where order_items.order_id = ?;
select * from order_payments where order_id = ?;

-- create sale
insert into orders (order_date, customer_id, discount_id, order_type, subtotal, discount_amount, tax_amount, total)
values (?, ?, ?, ?, ?, ?, ?, ?);
insert into order_items (order_id, product_id, discount_id, quantity, price, discount_amount)
values (?, ?, ?, ?, ?, ?);
insert into order_payments (order_id, payment_type, amount, payment_info)
values (?, ?, ?, ?);

-- bulk load customers
load data infile '/path/to/customers.csv'
into table customers
fields terminated by ','
lines terminated by '\n'
(name, phone, email, marketing_opt_in);

-- daily sold items
select id, name, category, sum(total_quantity) as total_quantity, sum(total_sales) as total_sales
from item_summaries
where order_date = ?
group by id, name, category
order by category, name;

-- daily revenue
select order_type,sum(total) as total_revenue
from orders
where order_date >= ? and order_date <= ?
group by order_type;

-- customer sales
select customers.id, customers.name, 
  sum(orders.total) as total_sales, count(*) as total_orders
from customers
join orders on customers.id = orders.customer_id
where orders.order_date >= ? and orders.order_date <= ?
group by customers.id, customers.name
order by total_sales desc;

-- general sales report
select ro.title, ro.report_order, t.name, 
sum(t.order_count) as order_count,
sum(t.total_quantity) as quantity,
sum(t.total_sales) as total_sales
from item_summaries t 
inner join reporting_order ro on ro.order_type = 'general' and ro.category = t.category
group by ro.title, ro.report_order, t.name 
order by ro.report_order, t.name;

-- typed sales report
select ro.title, ro.report_order, t.name, 
sum(t.order_count) as order_count,
sum(t.total_quantity) as quantity,
sum(t.total_sales) as total_sales
from item_summaries t 
inner join reporting_order ro on ro.order_type = t.order_type and ro.category = t.category
group by ro.title, ro.report_order, t.name 
order by ro.report_order, t.name;
