-- CustomerSales
select c.id, c.name, sum(o.total) as total_sales, count(o.id) as total_orders
from customers c
inner join orders o on o.customer_id = c.id
where o.order_date >= ? /* start_date */ and o.order_date <= ? /* end_date */
group by c.id, c.name
order by total_sales desc;

-- DailySoldItems
select id, name, category, 
sum(total_quantity) as total_quantity, 
sum(total_sales) as total_sales
from item_summaries
where order_date = ?
group by id, name, category
order by category, name;

-- GeneralSales
select ro.title, ro.report_order, t.name, 
sum(t.order_count) as order_count,
sum(t.total_quantity) as quantity,
sum(t.total_sales) as total_sales
from item_summaries t 
inner join reporting_order ro on ro.order_type = 'general' and ro.category = t.category
where t.order_date >= ? and t.order_date <= ?
group by ro.title, ro.report_order, t.name 
order by ro.report_order, t.name;


-- WeeklyTypedSales
select dim_date.year, dim_date.WEEK_OF_YEAR /* week */, 
ro.title, ro.report_order, t.name, 
sum(t.total_quantity) as quantity,
sum(t.total_sales) as total_sales
from item_summaries t 
inner join reporting_order ro on ro.order_type = t.order_type and ro.category = t.category
inner join dim_date on dim_date.date = t.order_date
where t.order_date >= ? /* start_date */ and t.order_date <= ? /* end_date */
group by dim_date.year, dim_date.WEEK_OF_YEAR,ro.title, ro.report_order, t.name 
order by dim_date.year, dim_date.WEEK_OF_YEAR, ro.report_order, t.name;