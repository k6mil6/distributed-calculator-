CREATE TABLE IF NOT EXISTS timeouts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    timeouts_values JSONB NOT NULL
);