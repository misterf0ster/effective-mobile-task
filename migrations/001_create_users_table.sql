CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255),
    gender VARCHAR(10),
    age INT,
    nationality VARCHAR(255)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
