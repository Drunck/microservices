CREATE TABLE IF NOT EXISTS orders (
    id bigserial PRIMARY KEY,
    book_id int NOT NULL,
    email text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);