-- +goose Up

CREATE TABLE refresh_tokens(
    token VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT fk_refresh_tokens
        FOREIGN KEY (user_id)
        REFERENCES users (id)
);

-- +goose Down

DROP TABLE refresh_tokens;