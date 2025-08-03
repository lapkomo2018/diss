-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT REFERENCES chats(id),
    sender_id INT REFERENCES users(id),
    content TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS users;
