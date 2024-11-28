-- +goose Up
DROP TABLE IF EXISTS roles_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS users_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS telegram_users;
DROP TABLE IF EXISTS users;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id            UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    created_at    TIMESTAMP           NOT NULL,
    updated_at    TIMESTAMP
);

CREATE TABLE telegram_users
(
    id             UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id        UUID REFERENCES users (id) UNIQUE NOT NULL,
    telegram_id    BIGINT UNIQUE,
    telegram_login TEXT UNIQUE                       NOT NULL,
    chat_id        BIGINT UNIQUE,
    created_at     TIMESTAMP                         NOT NULL NOT NULL,
    updated_at     TIMESTAMP
);

CREATE TABLE tokens
(
    token      UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id    UUID      NOT NULL REFERENCES users (id),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE profiles
(
    id            UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id       UUID REFERENCES users (id) UNIQUE NOT NULL,
    first_name    VARCHAR(255)                      NOT NULL,
    last_name     VARCHAR(255)                      NOT NULL,
    patronymic    VARCHAR(255),
    date_of_birth DATE                              NOT NULL,
    email         VARCHAR(255) UNIQUE               NOT NULL,
    phone         VARCHAR(50)                       NOT NULL,
    address       TEXT                              NOT NULL,
    created_at    TIMESTAMP                         NOT NULL,
    updated_at    TIMESTAMP
);

CREATE TABLE roles
(
    id          UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name        VARCHAR(255) UNIQUE NOT NULL,
    description TEXT                NOT NULL,
    created_at  TIMESTAMP           NOT NULL,
    updated_at  TIMESTAMP
);

INSERT INTO roles (name, description, created_at)
VALUES ('sys-admin', 'Отвечает за общую администрирование системы и контроль доступа', NOW()),
       ('security-admin', 'Управляет настройками безопасности и правами доступа пользователей', NOW()),
       ('project-admin', 'Управляет настройками и рабочими процессами, связанными с проектами', NOW()),
       ('manager', 'Контролирует команды и процессы внутри организации', NOW()),
       ('analyst', 'Проводит анализ данных и предоставляет рекомендации для принятия решений', NOW()),
       ('business-analyst', 'Анализирует бизнес-требования и помогает в определении объема проектов', NOW()),
       ('operator', 'Осуществляет работу и контроль ежедневных функций системы', NOW()),
       ('moderator', 'Модерирует пользовательский контент на соответствие правилам', NOW()),
       ('user', 'Обычный пользователь с доступом к базовым функциям системы', NOW()),
       ('guest', 'Пользователь с ограниченным доступом, в основном для просмотра', NOW()),
       ('superuser', 'Пользователь с расширенными привилегиями для критических действий в системе', NOW()),
       ('supporter', 'Оказывает поддержку и помощь пользователям системы', NOW()),
       ('developer', 'Отвечает за разработку и поддержку функций системы', NOW()),
       ('tester', 'Тестирует новые функции и изменения для обеспечения качества и функциональности', NOW())
;

CREATE TABLE users_roles
(
    id      UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id),
    role_id UUID NOT NULL REFERENCES roles (id),
    UNIQUE (user_id, role_id)
);

CREATE TABLE permissions
(
    id         UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    route_path VARCHAR(255) NOT NULL,
    method     VARCHAR(10)  NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP,
    UNIQUE (route_path, method)
);

INSERT INTO permissions (route_path, method, created_at)
VALUES ('/v1/users-roles', 'GET', NOW()),
       ('/v1/users/{id}', 'GET', NOW()),
       ('/v1/users', 'GET', NOW()),
       ('/v1/accounts', 'POST', NOW()),
       ('/v1/payments', 'POST', NOW()),
       ('/v1/payments/confirmation', 'POST', NOW()),
       ('/v1/payments/{id}', 'GET', NOW()),
       ('/v1/payments/{id}', 'PATCH', NOW()),
       ('/v1/payments/{id}/transactions', 'PATCH', NOW()),
       ('/v1/credit-applications', 'POST', NOW()),
       ('/v1/credit-applications/confirmation', 'POST', NOW()),
       ('/v1/credit-applications/{id}', 'GET', NOW()),
       ('/v1/credit-applications/{id', 'PATCH', NOW()),
       ('/v1/credits/{id}', 'GET', NOW()),
       ('/v1/credits', 'GET', NOW()),
       ('/v1/credits/{id}/payment-schedule', 'GET', NOW()),
       ('/v1/log-reports', 'POST', NOW()),
       ('/v1/reports', 'GET', NOW())
;

CREATE TABLE roles_permissions
(
    id            UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    role_id       UUID REFERENCES roles (id)       NOT NULL,
    permission_id UUID REFERENCES permissions (id) NOT NULL,
    UNIQUE (role_id, permission_id)
);

INSERT INTO roles_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
         CROSS JOIN permissions p
WHERE r.name = 'sys-admin';

INSERT INTO roles_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
         JOIN permissions p ON (
    (p.route_path = '/v1/users/{id}' AND p.method IN ('GET')) OR
    (p.route_path = '/v1/accounts' AND p.method IN ('POST')) OR
    (p.route_path = '/v1/payments' AND p.method IN ('POST')) OR
    (p.route_path = '/v1/payments/confirmation' AND p.method IN ('POST')) OR
    (p.route_path = '/v1/payments/{id}' AND p.method IN ('GET')) OR
    (p.route_path = '/v1/payments/{id}' AND p.method IN ('PATCH')) OR
    (p.route_path = '/v1/credit-applications' AND p.method IN ('POST')) OR
    (p.route_path = '/v1/credit-applications/confirmation' AND p.method IN ('POST')) OR
    (p.route_path = '/v1/credit-applications/{id}' AND p.method IN ('GET')) OR
    (p.route_path = '/v1/credits/{id}' AND p.method IN ('GET')) OR
    (p.route_path = '/v1/credits' AND p.method IN ('GET')) OR
    (p.route_path = '/v1/credits/{id}/payment-schedule' AND p.method IN ('GET')) OR
    (p.route_path = '/v1/reports' AND p.method IN ('GET'))
    )
WHERE r.name = 'user';


-- +goose Down
DROP TABLE IF EXISTS roles_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS users_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS telegram_users;
DROP TABLE IF EXISTS users;