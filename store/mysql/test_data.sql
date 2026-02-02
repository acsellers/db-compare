-- sample products for our electronics store
INSERT INTO products (id, name, category, price)
VALUES 
(1, 'Smasnug Laptop 15 inch', 'pc', 1200.00),
(2, 'Smasnug Laptop 13 inch', 'pc', 1000.00),
(3, 'Smasnug Phone', 'phone', 800.00),
(4, 'Smasnug Tablet', 'phone', 600.00),
(5, 'Smasnug Watch', 'phone', 200.00),
(6, 'Smasnug Headphones', 'audio', 50.00),
(7, 'Smasnug Earbuds', 'audio', 150.00),
(8, 'Smasnug Speakers', 'audio', 200.00);

-- sample customers
INSERT INTO customers (id, name, phone, email, marketing_opt_in)
VALUES 
(1, 'Alice Smith', '555-0101', 'alice@example.com', true),
(2, 'Bob Johnson', '555-0102', 'bob@example.com', false),
(3, 'Charlie Brown', '555-0103', 'charlie@example.com', true),
(4, 'Diana Prince', '555-0104', 'diana@example.com', false),
(5, 'Ethan Hunt', '555-0105', 'ethan@example.com', true);

-- sample discounts
INSERT INTO discounts (id, name, category, discount_type, discount)
VALUES 
(1, '10% off', 'electronics', 'percentage', 0.10),
(2, '15% off', 'electronics', 'percentage', 0.15),
(3, '$10 off', 'electronics', 'fixed', 10.00),
(4, '$25 off', 'electronics', 'fixed', 25.00),
(5, '$50 off', 'electronics', 'fixed', 50.00);

TRUNCATE TABLE reporting_order;
INSERT INTO reporting_order (id, category, order_type, report_order, title)
VALUES 
(1, 'pc', 'members', 1, 'Member PC Items'),
(2, 'phone', 'members', 2, 'Member Phone & Wearable Items'),
(3, 'audio', 'members', 3, 'Member Audio Items'),
(4, 'taxes', 'members', 4, 'Member Tax Items'),
(5, 'discounts', 'members', 5, 'Member Discounts'),
(6, 'payments', 'members', 6, 'Member Payments'),
(7, 'pc', 'non-members', 7, 'Non-Member PC Items'),
(8, 'phone', 'non-members', 8, 'Non-Member Phone & Wearable Items'),
(9, 'audio', 'non-members', 9, 'Non-Member Audio Items'),
(10, 'taxes', 'non-members', 10, 'Non-Member Tax Items'),
(11, 'discounts', 'non-members', 11, 'Non-Member Discounts'),
(12, 'payments', 'non-members', 12, 'Non-Member Payments'),
(13, 'pc', 'general', 1, 'PC Items'),
(14, 'phone', 'general', 2, 'Phone & Wearable Items'),
(15, 'audio', 'general', 3, 'Audio Items'),
(16, 'taxes', 'general', 4, 'Tax Items'),
(17, 'discounts', 'general', 5, 'Discounts'),
(18, 'payments', 'general', 6, 'Payments');

INSERT INTO orders (id, order_date, customer_id, discount_id, order_type, subtotal, discount_amount, tax_amount, total)
VALUES 
(1, '2026-01-01', 1, NULL, 'members', 1200.00, 0, 108.00, 1308.00),
(2, '2026-01-01', null, 5, 'non-members', 600.00, 50.00, 55.00, 605.00),
(3, '2026-01-01', 2, null, 'members', 2000.00, 80.00, 172.80, 2092.80);

INSERT INTO order_items (id, order_id, product_id, discount_id, quantity, price, discount_amount)
VALUES 
(1, 1, 1, NULL, 1, 1200.00, 0),
(2, 2, 4, NULL, 1, 600.00, 0),
(3, 3, 3, 1, 1, 800.00, 80.00),
(4, 3, 4, NULL, 2, 600.00, 0.00);

INSERT INTO order_payments (id, order_id, payment_type, amount, payment_info)
VALUES
(1, 1, 'cash', 608.00, null),
(2, 1, 'giftcard', 700.00, '{"card_number": "1234-5678-9012-3456"}'),
(3, 2, 'credit', 605.00, '{"card_type": "visa", "last_four": "1234"}'),
(4, 3, 'credit', 2092.80, '{"card_type": "visa", "last_four": "1234"}');

INSERT INTO payment_names (payment_type, name)
VALUES
('cash', 'Cash'),
('giftcard', 'Gift Card'),
('credit', 'Credit Card');

