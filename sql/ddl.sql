-- created at, updated at an deleted at columns will be added
--  on gorm auto migrate
CREATE TABLE cars (
    id SERIAL PRIMARY KEY,
    wheel_drive INT NOT NULL,
    type VARCHAR(50) NOT NULL,
    seats INT NOT NULL,
    transmission VARCHAR(20) NOT NULL,
    manufacturer VARCHAR(255) NOT NULL,
    car_model VARCHAR(255) NOT NULL,
    year INT NOT NULL,
    stock INT NOT NULL,
    rate_per_day DECIMAL NOT NULL,
    UNIQUE(manufacturer, car_model)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    password VARCHAR(255) NOT NULL,
    deposit DECIMAL NOT NULL DEFAULT 0
);

CREATE TABLE rentals (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) NOT NULL,
    car_id INT REFERENCES cars(id) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_price DECIMAL NOT NULL
);

CREATE TABLE top_ups (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) NOT NULL,
    amount DECIMAL NOT NULL
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    -- polymorphic FK to either rentals or top_ups
    purchase_id INT NOT NULL,
    payment_url VARCHAR(255) NOT NULL,
    status VARCHAR(100) NOT NULL,
    payment_method VARCHAR(100),
    total_payment DECIMAL
);

-- Insert dummy data into cars table
INSERT INTO cars (wheel_drive, type, seats, transmission, manufacturer, car_model, year, stock, rate_per_day) VALUES
(4, 'SUV', 5, 'Automatic', 'Toyota', 'RAV4', 2022, 10, 50.00),
(2, 'Sedan', 5, 'Manual', 'Honda', 'Civic', 2021, 8, 40.00),
(4, 'Truck', 2, 'Automatic', 'Ford', 'F-150', 2023, 5, 70.00),
(4, 'SUV', 7, 'Automatic', 'Chevrolet', 'Tahoe', 2022, 6, 60.00),
(2, 'Coupe', 4, 'Manual', 'BMW', 'M4', 2021, 3, 80.00);

-- Insert dummy data into users table
INSERT INTO users (name, email, password, deposit) VALUES
('John Doe', 'john@example.com', 'password123', 200.00),
('Jane Smith', 'jane@example.com', 'securepass', 150.00),
('Mike Johnson', 'mike@example.com', 'mikepass', 300.00),
('Emily Davis', 'emily@example.com', 'emilysecure', 100.00),
('David Wilson', 'david@example.com', 'david123', 250.00);

-- Insert dummy data into rentals table
INSERT INTO rentals (user_id, car_id, start_date, end_date, total_price) VALUES
(1, 1, '2024-09-01', '2024-09-05', 200.00),
(2, 2, '2024-09-10', '2024-09-12', 80.00),
(3, 3, '2024-09-15', '2024-09-20', 350.00),
(4, 4, '2024-09-05', '2024-09-08', 180.00),
(5, 5, '2024-09-20', '2024-09-25', 400.00);

-- Insert dummy data into top_ups table
INSERT INTO top_ups (user_id, amount) VALUES
(1, 50.00),
(2, 30.00),
(3, 70.00),
(4, 100.00),
(5, 25.00);

-- Insert dummy data into payments table
INSERT INTO payments (purchase_id, payment_url, status, payment_method, total_payment) VALUES
(1, 'http://payment.com/1', 'Completed', 'Credit Card', 200.00),
(2, 'http://payment.com/2', 'Completed', 'PayPal', 80.00),
(3, 'http://payment.com/3', 'Pending', 'Credit Card', 350.00),
(4, 'http://payment.com/4', 'Completed', 'Debit Card', 180.00),
(5, 'http://payment.com/5', 'Failed', 'Credit Card', 400.00);