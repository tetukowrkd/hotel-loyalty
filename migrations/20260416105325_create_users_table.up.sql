CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,

    role_id UUID,

    phone TEXT UNIQUE,

    is_active BOOLEAN DEFAULT TRUE,
    is_deleted BOOLEAN DEFAULT FALSE,

    email_verified_at TIMESTAMP,
    last_login_at TIMESTAMP,
    login_attempt INT DEFAULT 0,

    reset_password_token TEXT,
    reset_password_expired_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_role
        FOREIGN KEY(role_id)
        REFERENCES roles(id)
        ON DELETE SET NULL
);