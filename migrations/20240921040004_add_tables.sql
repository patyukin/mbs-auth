-- +goose Up
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
    user_id        UUID REFERENCES users (id) ON UPDATE CASCADE ON DELETE NO ACTION UNIQUE NOT NULL,
    telegram_id    BIGINT UNIQUE,
    telegram_login TEXT UNIQUE                                                             NOT NULL,
    chat_id        BIGINT UNIQUE,
    created_at     TIMESTAMP                                                               NOT NULL NOT NULL,
    updated_at     TIMESTAMP
);

CREATE TABLE tokens
(
    token      UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id    UUID      NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE NO ACTION,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE profiles
(
    id            UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id       UUID REFERENCES users (id) ON UPDATE CASCADE ON DELETE NO ACTION UNIQUE NOT NULL,
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

CREATE TABLE roles
(
    id          UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name        VARCHAR(255) UNIQUE NOT NULL,
    description TEXT                NOT NULL,
    created_at  TIMESTAMP           NOT NULL,
    updated_at  TIMESTAMP
);

INSERT INTO roles (name, description, created_at)
VALUES ('system-admin', 'Отвечает за общую администрирование системы и контроль доступа', NOW()),
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
    id         UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id    UUID REFERENCES users (id) ON UPDATE CASCADE ON DELETE NO ACTION,
    role_id    UUID REFERENCES roles (id) ON UPDATE CASCADE ON DELETE NO ACTION,
    created_at TIMESTAMP NOT NULL
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
VALUES ('/auth/users-roles', 'GET', NOW()),
       ('/auth/users-roles', 'POST', NOW()),
       ('/auth/users-roles/{id}', 'GET', NOW()),
       ('/auth/users-roles/{id}', 'PUT', NOW()),
       ('/auth/users-roles/{id}', 'DELETE', NOW()),
       ('/payments/payments', 'GET', NOW()),
       ('/payments/payments', 'POST', NOW()),
       ('/payments/payments/{id}', 'GET', NOW()),
       ('/payments/payments/{id}', 'PUT', NOW()),
       ('/payments/payments/{id}', 'DELETE', NOW()),
       ('/payments/accounts', 'GET', NOW()),
       ('/payments/accounts', 'POST', NOW()),
       ('/payments/accounts/{id}', 'GET', NOW()),
       ('/payments/accounts/{id}', 'PUT', NOW()),
       ('/payments/accounts/{id}', 'DELETE', NOW()),
       ('/payments/transactions', 'GET', NOW()),
       ('/payments/transactions/{id}', 'GET', NOW()),
       ('/reports/reports', 'GET', NOW()),
       ('/reports/reports', 'POST', NOW()),
       ('/reports/reports/{id}', 'GET', NOW()),
       ('/reports/reports/{id}', 'PUT', NOW()),
       ('/reports/reports/{id}', 'DELETE', NOW())
;

CREATE TABLE roles_permissions
(
    id            UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    role_id       UUID REFERENCES roles (id) ON UPDATE CASCADE ON DELETE NO ACTION,
    permission_id UUID REFERENCES permissions (id) ON UPDATE CASCADE ON DELETE NO ACTION,
    created_at    TIMESTAMP NOT NULL
);

INSERT INTO roles_permissions (role_id, permission_id, created_at)
VALUES ((SELECT id FROM roles WHERE name = 'system-admin'),
        (SELECT id FROM permissions WHERE route_path = '/auth/users-roles' AND method = 'GET'), NOW()),
       ((SELECT id FROM roles WHERE name = 'system-admin'),
        (SELECT id FROM permissions WHERE route_path = '/auth/users-roles' AND method = 'POST'), NOW()),
       ((SELECT id FROM roles WHERE name = 'system-admin'),
        (SELECT id FROM permissions WHERE route_path = '/auth/users-roles/{id}' AND method = 'GET'), NOW()),
       ((SELECT id FROM roles WHERE name = 'system-admin'),
        (SELECT id FROM permissions WHERE route_path = '/auth/users-roles/{id}' AND method = 'PUT'), NOW()),
       ((SELECT id FROM roles WHERE name = 'system-admin'),
        (SELECT id FROM permissions WHERE route_path = '/auth/users-roles/{id}' AND method = 'DELETE'), NOW()),
       ((SELECT id FROM roles WHERE name = 'manager'),
        (SELECT id FROM permissions WHERE route_path = '/auth/users-roles' AND method = 'GET'), NOW()),
       ((SELECT id FROM roles WHERE name = 'analyst'),
        (SELECT id FROM permissions WHERE route_path = '/auth/users-roles' AND method = 'GET'), NOW()),
       ((SELECT id FROM roles WHERE name = 'business-analyst'),
        (SELECT id FROM permissions WHERE route_path = '/auth/users-roles' AND method = 'GET'), NOW()),

       ((SELECT id FROM roles WHERE name = 'superuser'),
        (SELECT id FROM permissions WHERE route_path = '/payments/payments' AND method = 'GET'), NOW()),
       ((SELECT id FROM roles WHERE name = 'superuser'),
        (SELECT id FROM permissions WHERE route_path = '/payments/payments' AND method = 'POST'), NOW()),
       ((SELECT id FROM roles WHERE name = 'superuser'),
        (SELECT id FROM permissions WHERE route_path = '/payments/payments/{id}' AND method = 'POST'), NOW()),
       ((SELECT id FROM roles WHERE name = 'superuser'),
        (SELECT id FROM permissions WHERE route_path = '/payments/payments/{id}' AND method = 'PUT'), NOW()),
       ((SELECT id FROM roles WHERE name = 'superuser'),
        (SELECT id FROM permissions WHERE route_path = '/payments/payments/{id}' AND method = 'DELETE'), NOW()),
       ((SELECT id FROM roles WHERE name = 'user'),
        (SELECT id FROM permissions WHERE route_path = '/payments/payments' AND method = 'GET'), NOW()),
       ((SELECT id FROM roles WHERE name = 'user'),
        (SELECT id FROM permissions WHERE route_path = '/payments/payments' AND method = 'POST'), NOW()),
       ((SELECT id FROM roles WHERE name = 'user'),
        (SELECT id FROM permissions WHERE route_path = '/payments/payments/{id}' AND method = 'GET'), NOW()),
       ((SELECT id FROM roles WHERE name = 'user'),
        (SELECT id FROM permissions WHERE route_path = '/payments/payments/{id}' AND method = 'PUT'), NOW()),
       ((SELECT id FROM roles WHERE name = 'user'),
        (SELECT id FROM permissions WHERE route_path = '/payments/payments/{id}' AND method = 'DELETE'), NOW())
;


-- +goose Down
DROP TABLE IF EXISTS roles_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS users_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS telegram_users;
DROP TABLE IF EXISTS users;