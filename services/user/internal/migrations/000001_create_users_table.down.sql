CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);