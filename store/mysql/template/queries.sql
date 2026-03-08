-- 01 - Get Sale

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

-- name: GetLocationSales
select orders.*, coalesce(customers.name, 'Guest') as customer_name, locations.name as location_name
from orders 
left join customers on orders.customer_id = customers.id
left join locations on locations.id = orders.location_id
where orders.location_id = ? and orders.order_date >= ? and orders.order_date <= ?;

-- name: GetSalesItems
select order_items.*, 
  products.name, products.category 
from order_items 
left join products on order_items.product_id = products.id
where order_items.order_id IN (?);

-- name: GetSalesPayments
select * from order_payments where order_id IN (?);


-- 02 - Create Sale

-- name: CreateSale
insert into orders (order_date, customer_id, location_id, discount_id, order_type, subtotal, discount_amount, tax_amount, total)
values (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: CreateSaleItems
insert into order_items (order_id, product_id, discount_id, quantity, price, discount_amount)
values (?, ?, ?, ?, ?, ?);

-- name: CreateSalePayments
insert into order_payments (order_id, payment_type, amount, payment_info)
values (?, ?, ?, ?);

-- 03 - Sale Search

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
and (orders.location_id = @location_id or @location_id = 0)
order by orders.order_date desc;

-- 04 - Bulk Load Customers

-- name: InsertCustomers
INSERT INTO customers (name, email, phone, join_location_id) VALUES (?, ?, ?, ?);

-- name: InsertCustomersBulk
INSERT INTO customers  (name, phone, email, join_location_id)  
VALUES (?, ?, ?, ?);

-- 05 - Customer Update

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

-- 06 - JSON




-- 07 - With Queries

-- name: WithProducts
WITH src as (
  select o.location_id, date(o.order_date), p.name, p.price, count(distinct o.id) as order_count, 
  sum(oi.item_total) as total, sum(oi.quantity) as quantity, 
  sum(oi.discount_amount) as alt_total
  from orders o
  inner join order_items oi on oi.order_id = o.id
  inner join products p on oi.product_id = p.id
  where o.order_date >= '2024-01-01' and o.order_date < '2025-01-01'
  and p.category = ?
  group by p.name, p.price, o.location_id , date(o.order_date)
)
-- name: WithSalesTax
With src as (
  select l.id, date(o.order_date), concat(l.name, ' Sales Tax'), l.tax_rate, count(distinct o.id) as order_count,
  sum(o.tax_amount ) as total, count(distinct o.id) as quantity,
  sum(o.subtotal ) as alt_total
  from orders o
  inner join locations l  on l.id = o.location_id 
  where o.order_date >= '2024-01-01' and o.order_date < '2025-01-01'
  group by l.name, l.tax_rate, l.id, date(o.order_date)
)
-- name: WithPayments
With src as (
  SELECT o.location_id, date(o.order_date), case when op.payment_type = 'cash' then 'Cash'
  when op.payment_type = 'credit_card' then op.payment_info->>'$.card_type'
  when op.payment_type = 'purchase_order' then 'Purchase Order'
  when op.payment_type = 'gift_card' then 'Gift Card'
  else op.payment_type end as group_value, 0, count(distinct o.id), sum(op.amount ) as total,
  count(op.id) as quantity, sum(o.total) as alt_total
  from orders o 
  inner join order_payments op on op.order_id  = o.id 
  where o.order_date >= '2024-01-01' and o.order_date < '2025-01-01'
  group by o.location_id, date(o.order_date), group_value
)
-- name: WithDiscounts
With src as (
  select coalesce(o.location_id, o2.location_id), coalesce(date(o.order_date), date(o2.order_date)), d.name, d.discount, 
  count(o.id) + count(distinct o2.id) as order_count,
  sum(o.discount_amount) + sum(oi.discount_amount) as total,
  count(o.id) + sum(oi.quantity ) as quantity,
  sum(o.subtotal ) + sum(oi.item_total ) as alt_total
  from discounts d 
  left join orders o on d.id = o.discount_id 
  left join order_items oi on oi.discount_id = d.id
  left join orders o2 on o2.id = oi.order_id 
  where (o.order_date >= '2024-01-01' and o.order_date < '2025-01-01') or 
  (o2.order_date >= '2024-01-01' and o2.order_date < '2025-01-01')
  group by d.name, d.discount, coalesce(o.location_id, o2.location_id), coalesce(date(o.order_date), date(o2.order_date))
)
-- name: WithTotals
With src as (
  select o.location_id, date(o.order_date), o.order_type , 0, count(o.id) as order_count, 
  sum(o.total ) as total, count(o.id) as quantity, 
  sum(o.discount_amount) as alt_total
  from orders o
  where o.order_date >= '2024-01-01' and o.order_date < '2025-01-01'
  group by o.order_type, o.location_id , date(o.order_date) 
)

;

-- 08 - Basic Grouping

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

-- name: DailySoldItems
select id, name, category, 
cast(sum(total_quantity) as signed) as total_quantity, 
cast(sum(total_sales) as double) as total_sales
from item_summaries
where order_date = ?
group by id, name, category
order by category, name;

-- 09 - Advanced Grouping

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