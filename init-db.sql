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
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- Внешний ключ для связи с таблицей products
    CONSTRAINT fk_cart_product 
        FOREIGN KEY (product_id) 
        REFERENCES products(id) 
        ON DELETE CASCADE
);

-- Создание таблицы favorites
CREATE TABLE favorites (
    id SERIAL PRIMARY KEY,
    user_email VARCHAR(50) NOT NULL,
    product_id INTEGER NOT NULL,
    -- Внешний ключ для связи с таблицей products
    CONSTRAINT fk_favorites_product 
        FOREIGN KEY (product_id) 
        REFERENCES products(id) 
        ON DELETE CASCADE,
    -- Уникальность комбинации user_email и product_id
    CONSTRAINT unique_user_product 
        UNIQUE (user_email, product_id)
);

-- Создание таблицы user_sessions
CREATE TABLE user_sessions (
    token VARCHAR(64) PRIMARY KEY,
    user_email VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание индексов для улучшения производительности
-- Индекс для cart по user_email (часто используется для поиска корзины пользователя)
CREATE INDEX idx_cart_user_email ON cart(user_email);

-- Индекс для cart по product_id
CREATE INDEX idx_cart_product_id ON cart(product_id);

-- Индекс для favorites по user_email
CREATE INDEX idx_favorites_user_email ON favorites(user_email);

-- Индекс для favorites по product_id
CREATE INDEX idx_favorites_product_id ON favorites(product_id);

-- Индекс для product_images по product_id
CREATE INDEX idx_product_images_product_id ON product_images(product_id);

-- Индекс для user_sessions по user_email
CREATE INDEX idx_user_sessions_user_email ON user_sessions(user_email);

-- Индекс для user_sessions по created_at (очистка старых сессий)
CREATE INDEX idx_user_sessions_created_at ON user_sessions(created_at);

-- Вставка данных в таблицу products
INSERT INTO products (id, name, description, price, imageurl, image_urls) VALUES
(2, 'Конвектор серый', 'Описание конвектора', 2000.00, 'RedHeart/img.png', ARRAY['RedHeart/img.png', 'RedHeart/img_1.png']),
(4, 'Серебристый конструктор', 'Описание конструктора', 120000.00, 'MoonSergi/img.png', ARRAY['MoonSergi/img.png', 'MoonSergi/img_1.png']);