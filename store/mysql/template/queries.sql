-- name: GetDiscount
SELECT * FROM discounts WHERE id IN (?);

-- name: GetProducts
SELECT * FROM products WHERE id IN (?);

-- name: CustomerExists
select count(*) from customers where id = ?;

-- name: GetSale
select orders.*, coalesce(customers.name, 'Guest') as customer_name, locations.name as location_name
from orders 
left join customers on orders.customer_id = customers.id
left join locations on locations.id = orders.location_id
where orders.id = ?;

-- name: GetSaleItems
select order_items.*, 
  products.name, products.category 
from order_items 
left join products on order_items.product_id = products.id
where order_items.order_id = ?;

-- name: GetSalePayments
select * from order_payments where order_id = ?;

-- name: CreateSale
insert into orders (order_date, customer_id, location_id, discount_id, order_type, subtotal, discount_amount, tax_amount, total)
values (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: CreateSaleItems
insert into order_items (order_id, product_id, discount_id, quantity, price, discount_amount)
values (?, ?, ?, ?, ?, ?);

-- name: CreateSalePayments
insert into order_payments (order_id, payment_type, amount, payment_info)
values (?, ?, ?, ?);

-- name: InsertCustomers
INSERT INTO customers (name, email, phone, join_location_id) VALUES (?, ?, ?, ?);

-- name: InsertCustomersBulk
INSERT INTO customers  (name, phone, email, join_location_id)  
VALUES (?, ?, ?, ?);

-- name: UpdateCustomerByExternalID
update customers set name = ?, email = ?, phone = ? where external_id = ?;

-- name: CreateCustomerTempTable
CREATE TEMPORARY TABLE customer_temp (
  external_id VARCHAR(12),
  name VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(24),
);

-- name: InsertCustomersTemp
INSERT INTO customer_temp (external_id, name, email, phone) VALUES (?, ?, ?, ?);

-- name: UpdateCustomerNamesFromTemp
update customers c join customer_temp ct on c.external_id = ct.external_id
set c.name = ct.name;

-- name: UpdateCustomerEmailsFromTemp
update customers c join customer_temp ct on c.external_id = ct.external_id
set c.email = ct.email;

-- name: UpdateCustomerPhonesFromTemp
update customers c join customer_temp ct on c.external_id = ct.external_id
set c.phone = ct.phone;

-- name: DailySoldItems
select id, name, category, 
cast(sum(total_quantity) as signed) as total_quantity, 
cast(sum(total_sales) as double) as total_sales
from item_summaries
where order_date = ?
group by id, name, category
order by category, name;

-- name: DailyRevenue
select order_type,sum(total) as total_revenue
from orders
where order_date >= sqlc.arg(start_date) and order_date <= sqlc.arg(end_date)
group by order_type;

-- name: CustomerSales
select customers.id, customers.name, 
  sum(orders.total) as total_sales, count(*) as total_orders
from customers
join orders on customers.id = orders.customer_id
where orders.order_date >= sqlc.arg(start_date) and orders.order_date <= sqlc.arg(end_date)
group by customers.id, customers.name
order by total_sales desc;

-- name: GeneralSales
select ro.title, ro.report_order, t.name, 
sum(t.order_count) as order_count,
sum(t.total_quantity) as quantity,
sum(t.total_sales) as total_sales
from item_summaries t 
inner join reporting_order ro on ro.order_type = 'general' and ro.category = t.category
where t.order_date >= sqlc.arg(start_date) and t.order_date <= sqlc.arg(end_date)
and (t.location_id = sqlc.arg(location_id) or sqlc.arg(location_id) = 0)
group by ro.title, ro.report_order, t.name 
order by ro.report_order, t.name;

-- name: WeeklyTypedSales
select dim_date.year, dim_date.WEEK_OF_YEAR, ro.title, ro.report_order, t.name, 
sum(t.order_count) as order_count,
sum(t.total_quantity) as quantity,
sum(t.total_sales) as total_sales
from item_summaries t 
inner join reporting_order ro on ro.order_type = t.order_type and ro.category = t.category
inner join dim_date on dim_date.date = t.order_date
where t.order_date >= sqlc.arg(start_date) and t.order_date <= sqlc.arg(end_date)
and (t.location_id = sqlc.arg(location_id) or sqlc.arg(location_id) = 0)
group by dim_date.year, dim_date.WEEK_OF_YEAR,ro.title, ro.report_order, t.name 
order by dim_date.year, dim_date.WEEK_OF_YEAR, ro.report_order, t.name;

-- name: SearchSales
select orders.*, coalesce(customers.name, 'Guest') as customer_name, locations.name as location_name
from orders 
left join customers on orders.customer_id = customers.id
left join locations on locations.id = orders.location_id
where (customers.name like @name or @name = '')
and (orders.order_date >= @order_start or @order_start = '')
and (orders.order_date <= @order_end or @order_end = '')
and (orders.order_type = @order_type or @order_type = '')
and (orders.total >= @total_start or @total_start = 0)
and (orders.total <= @total_end or @total_end = 0)
order by orders.order_date desc;

-- name: WithProducts
WITH src as (
  select 
)
-- name: WithSalesTax

-- name: WithPayments

-- name: WithDiscounts

-- name: WithTotals
