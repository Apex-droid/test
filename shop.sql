CREATE TABLE rack (
    rack_id SERIAL PRIMARY KEY,
    rack_name VARCHAR(50)
);
CREATE TABLE product (
    product_id SERIAL PRIMARY KEY,
    product_name VARCHAR(50),
    rack_prime_id INT,
    FOREIGN KEY (rack_prime_id) REFERENCES rack(rack_id)
);

CREATE TABLE product_racks (
    rack_id INT,
    product_id INT,
    PRIMARY KEY (rack_id, product_id),
    FOREIGN KEY (rack_id) REFERENCES rack(rack_id),
    FOREIGN KEY (product_id) REFERENCES product(product_id)
);

CREATE TABLE orders (
  order_id SERIAL PRIMARY KEY,
  order_number INT
);

CREATE TABLE order_details (
    order_id INT,
    product_id INT,
    quantity INT,
    FOREIGN KEY (order_id) REFERENCES orders(order_id),
    FOREIGN KEY (product_id) REFERENCES product(product_id)
);
INSERT INTO rack (rack_name) VALUES
('А'),
('Б'),
('В'),
('З'),
('Ж');

INSERT INTO orders (order_number) VALUES 
(10),
(11),
(14),
(15);

INSERT INTO product (product_id, product_name, rack_prime_id) VALUES
(1,'Ноутбук', 1),
(2,'Телевизор', 1),
(3,'Телефон', 2),
(4, 'Системный блок',5),
(5, 'Часы',5),
(6, 'Микрофон',5);


INSERT INTO product_racks (product_id, rack_id) VALUES
(3, 3),
(3, 4),
(5, 1);

INSERT INTO order_details (order_id, product_id, quantity) VALUES
(1, 1, 2),
(1, 3, 1),
(1, 6, 1),
(2, 2, 3),
(3, 4, 4),
(3, 1, 3),
(4, 5, 1);
