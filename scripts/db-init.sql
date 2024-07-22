CREATE TABLE item_table (
    name VARCHAR(64) PRIMARY KEY,
    bar_code VARCHAR(32),
    country_of_origin VARCHAR(64),
    tot_quantity INTEGER,
    available BOOLEAN,
    expiry_date DATE,
    price INTEGER

);

INSERT INTO item_table (name, bar_code, country_of_origin, tot_quantity, available, expiry_date, price)
VALUES ('apples', 'A1B2C3D4E5F6G7H8', 'UK', 300, true, '2024-07-16', 1),
       ('herring', 'A1B2C3D4E5F6G7L8', 'Norway', 300, true, '2024-07-17', 4),
       ('meatballs', 'A1B2C3D4E5F6G7G8', 'Sweden', 300, true, '2024-07-18', 3),
       ('potatoes', 'A1B2C3D4E5F6G7K8', 'Ireland', 300, true, '2024-07-19', 2),
       ('brie', 'A1B2C3D4E5F6G7F8', 'France', 300, true, '2024-07-20', 5);