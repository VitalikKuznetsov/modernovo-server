# 🛒 Modernovo

Modernovo — полнофункциональный интернет-магазин с RESTful API бэкендом на Go, использующий PostgreSQL и Docker.  
Проект реализует пользовательскую систему, корзину покупок, избранные товары, административную панель и дополнительные сервисы.

## 🚀 Возможности

### 👤 Пользовательская система
- Регистрация и аутентификация пользователей
- Управление профилем (имя, телефон)
- Система сессий с токенами
- Избранные товары
- Корзина покупок

### 🛍️ Функционал магазина
- Просмотр каталога товаров с пагинацией
- Детальная информация о товарах
- Поиск товаров
- Добавление и удаление товаров из избранного
- Управление корзиной (добавление, изменение количества, удаление)

### 👨‍💼 Административная панель
- CRUD операции с товарами
- Просмотр всех товаров с полной информацией
- Создание, обновление и удаление товаров

### 🌤️ Дополнительные функции
- Погодный виджет (интеграция с OpenWeatherMap)
- CORS поддержка
- Автоматическая инициализация базы данных

## 🏗️ Технологический стек
- Go 1.21
- PostgreSQL 17
- Docker, Docker Compose
- Gorilla Mux
- RESTful JSON API
- Go testing


## 🚀 Быстрый старт

### Предварительные требования
- Docker
- Docker Compose
- Go 1.21+ (для локальной разработки)

### Запуск через Docker

```bash
git clone <repository-url>
cd modernovo
docker-compose up --build
```

## 🔧 API Endpoints

### Аутентификация

POST /api/register

POST /api/login

POST /api/logout

GET /api/profile

PUT /api/profile

### Товары

GET /api/products

GET /api/products/{id}

GET /api/products/{id}/detail

### Избранное

GET /api/favorites

POST /api/favorites

DELETE /api/favorites

GET /api/favorites/check

### Корзина

GET /api/cart

POST /api/cart

PUT /api/cart

DELETE /api/cart

POST /api/cart/clear

### Административные функции

GET /api/admin/products

POST /api/admin/products

PUT /api/admin/products/{id}

DELETE /api/admin/products/{id}

### Дополнительные

GET /api/weather

## 🧪 Тестирование

```bash
go test -v
```
