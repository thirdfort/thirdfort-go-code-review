-- Active: 1702997645398@@127.0.0.1@5432@consumer-api@public
CREATE TABLE IF NOT EXISTS recommendation (
	created_at timestamptz NOT NULL,
	updated_at timestamptz,
	deleted_at timestamptz,
	id text PRIMARY KEY,
    transaction_id text,
    status text,
    slug text,
    task text,
    attributes jsonb,
    CONSTRAINT "fk-recommendation"
      FOREIGN KEY (transaction_id)
      REFERENCES transaction(id)
      ON DELETE CASCADE
);
