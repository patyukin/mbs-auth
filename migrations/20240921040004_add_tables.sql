-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id            UUID                         DEFAULT uuid_generate_v4() PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    role          VARCHAR(20)         NOT NULL DEFAULT 'user',
    created_at    TIMESTAMP           NOT NULL,
    updated_at    TIMESTAMP
);

CREATE TABLE telegram_users
(
    id             UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id      UUID REFERENCES users (id) ON UPDATE CASCADE ON DELETE NO ACTION UNIQUE NOT NULL,
    telegram_login TEXT                                                                    NOT NULL,
    telegram_id    BIGINT UNIQUE,
    chat_id        BIGINT UNIQUE,
    created_at     TIMESTAMP                                                               NOT NULL NOT NULL,
    updated_at     TIMESTAMP
);

CREATE TABLE tokens
(
    token      UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id  UUID      NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE NO ACTION,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL
);

CREATE TABLE profiles
(
    id            UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id     UUID REFERENCES users (id) ON UPDATE CASCADE ON DELETE NO ACTION UNIQUE NOT NULL,
    first_name    VARCHAR(255)                                                            NOT NULL,
    last_name     VARCHAR(255)                                                            NOT NULL,
    patronymic    VARCHAR(255),
    date_of_birth DATE                                                                    NOT NULL,
    email         VARCHAR(255) UNIQUE                                                     NOT NULL,
    phone         VARCHAR(50)                                                             NOT NULL,
    address       TEXT                                                                    NOT NULL,
    created_at    TIMESTAMP                                                               NOT NULL,
    updated_at    TIMESTAMP
);


-- +goose Down
