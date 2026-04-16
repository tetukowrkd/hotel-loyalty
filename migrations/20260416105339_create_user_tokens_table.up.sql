CREATE TABLE user_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    refresh_token TEXT UNIQUE NOT NULL,
    user_agent TEXT,
    ip_address VARCHAR(50),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);