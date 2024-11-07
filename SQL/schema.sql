use mmackay;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS products;

CREATE TABLE customers (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    FirstName varchar(255),
    LastName varchar(255),
    email varchar(255)
);

CREATE TABLE products (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    product_name varchar(255),
    image_name varchar(255),
    price decimal(6, 2),
    in_stock int
);

CREATE TABLE orders (
    ID int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    product_id int,
    customer_id int,
    quantity int,
    price decimal(6, 2),
    tax decimal(6, 2),
    donation decimal(6, 2),
    timestamp bigint,
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);

INSERT INTO customers (FirstName, LastName, email) VALUES ('Mickey', 'Mouse', 'mmouse@mines.edu');
INSERT INTO customers (FirstName, LastName, email) VALUES ('Mallorie', 'Mackay', 'mmackay@mines.edu');

INSERT INTO products (product_name, image_name, price, in_stock) VALUES ('Lily Lamp', 'Lily.jpg', 199.99, 10);
INSERT INTO products (product_name, image_name, price, in_stock) VALUES ('Lotus Lamp', 'Lotus.jpg', 2212.50, 0);
INSERT INTO products (product_name, image_name, price, in_stock) VALUES ('Rose Lamp', 'Rose.jpg', 95.29, 3);
