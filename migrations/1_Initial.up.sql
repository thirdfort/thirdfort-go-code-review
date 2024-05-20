-- Actors is used for matching transactions by id
CREATE TABLE IF NOT EXISTS actor (
  id text PRIMARY KEY,
  name text,
  mobile text,
  email text,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS fingerprint (
  fingerprint text PRIMARY KEY,
  actor_id text NOT NULL,
  pa_actor_id text,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  CONSTRAINT "fk-actor"
    FOREIGN KEY (actor_id)
    REFERENCES actor(id)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS transaction (
  id text PRIMARY KEY,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  name text,
  tenant_name text,
  actor_name text,
  ref text,
  status text,
  event_id text,
  actor_id text
);

CREATE TABLE IF NOT EXISTS dob (
  id text PRIMARY KEY,
  transaction_id text NOT NULL,
  type text NOT NULL,
  pa_type text NOT NULL,
  status text NOT NULL,
  reason_code text,
  dob timestamptz,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  expectation_id text NOT NULL,
  CONSTRAINT "fk-transaction"
    FOREIGN KEY (transaction_id)
    REFERENCES transaction(id)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS name (
  id text PRIMARY KEY,
  transaction_id text NOT NULL,
  type text NOT NULL,
  pa_type text NOT NULL,
  status text NOT NULL,
  reason_code text,
  first text,
  last text,
  other text,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  expectation_id text NOT NULL,
  CONSTRAINT "fk-transaction"
    FOREIGN KEY (transaction_id)
    REFERENCES transaction(id)
    ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS address (
  id text PRIMARY KEY,
  transaction_id text NOT NULL,
  type text NOT NULL,
  pa_type text NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  status text NOT NULL,
  reason_code text,
  address_1 text,
  address_2 text,
  building_name text,
  building_number text,
  country text,
  flat_number text,
  postcode text,
  state text,
  street text,
  sub_street text,
  town text,
  expectation_id text NOT NULL,
  CONSTRAINT "fk-transaction"
    FOREIGN KEY (transaction_id)
    REFERENCES transaction(id)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS document (
  id text PRIMARY KEY,
  transaction_id text NOT NULL,
  type TEXT NOT NULL,
  pa_type text NOT NULL,
  status text NOT NULL,
  reason_code text,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  document_id text NOT NULL,
  expectation_id text NOT NULL,
  CONSTRAINT "fk-transaction"
    FOREIGN KEY (transaction_id)
    REFERENCES transaction(id)
    ON DELETE CASCADE
);

-- Expectation is a catch-all table for all tasks that don't contain data
CREATE TABLE IF NOT EXISTS expectation (
  id text PRIMARY KEY,
  transaction_id text NOT NULL,
  type text NOT NULL,
  pa_type text NOT NULL,
  status text NOT NULL,
  reason_code text,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  expectation_id text NOT NULL,
  CONSTRAINT "fk-transaction"
    FOREIGN KEY (transaction_id)
    REFERENCES transaction(id)
    ON DELETE CASCADE
);
