CREATE TABLE users (
    email VARCHAR(50) NOT NULL,
    password VARCHAR(30) NOT NULL,
    name VARCHAR(50) NOT NULL,
    phonenumber VARCHAR(20) NOT NULL
);

-- Создание таблицы products
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    imageurl VARCHAR(255),
    image_urls TEXT[]
);

-- Создание таблицы cart
CREATE TABLE cart (
    id SERIAL PRIMARY KEY,
    user_email VARCHAR(255) NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER DEFAULT 1,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы favorites
CREATE TABLE favorites (
    id SERIAL PRIMARY KEY,
    user_email VARCHAR(50) NOT NULL,
    product_id INTEGER NOT NULL
);

-- Создание таблицы user_sessions
CREATE TABLE user_sessions (
    token VARCHAR(64) PRIMARY KEY,
    user_email VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO products (id, name, description, price, imageurl, image_urls) VALUES
(2, 'Конвектор серый', 'Описание конвектора', 2000.00, 'RedHeart/img.png', ARRAY['RedHeart/img.png', 'RedHeart/img_1.png']),
(4, 'Серебристый конструктор', 'Описание конструктора', 120000.00, 'MoonSergi/img.png', ARRAY['MoonSergi/img.png', 'MoonSergi/img_1.png']);