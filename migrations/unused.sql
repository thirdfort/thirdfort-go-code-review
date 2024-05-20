
CREATE TABLE IF NOT EXISTS sof_purchaser (
  id text PRIMARY KEY,
  transaction_id text NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  property_expectation_id text,
  questionnaire_sof_expectation_id text,
  PRIMARY KEY (transaction_id, task_id),
  CONSTRAINT "fk-transaction"
    FOREIGN KEY (transaction_id)
    REFERENCES transaction(id)
    ON DELETE CASCADE
);

-- sof_giftor
CREATE TABLE IF NOT EXISTS questionnaire (
  id text PRIMARY KEY,
  transaction_id text NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  expectation_id text,
  PRIMARY KEY (transaction_id, task_id),
  CONSTRAINT "fk-transaction"
    FOREIGN KEY (transaction_id)
    REFERENCES transaction(id)
    ON DELETE CASCADE
);

-- sof_giftor
CREATE TABLE IF NOT EXISTS gift (
  id text PRIMARY KEY,
  transaction_id text NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  expectation_id text,
  PRIMARY KEY (transaction_id, task_id),
  CONSTRAINT "fk-transaction"
    FOREIGN KEY (transaction_id)
    REFERENCES transaction(id)
    ON DELETE CASCADE
);

-- ----------------

CREATE TABLE IF NOT EXISTS expectation (
  id text PRIMARY KEY,
  transaction_id text NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  pa_id text NOT NULL,
  status text,
  expectation text NOT NULL,
  data text,
  CONSTRAINT "fk-transaction"
    FOREIGN KEY (transaction_id)
    REFERENCES transaction(id)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS parts (
  id text PRIMARY KEY,
  expectation_id text NOT NULL,
  status text,
  type text NOT NULL,
  created_at timestamptz NOT NULL,
  deleted_at timestamptz,
  updated_at timestamptz,
  data text
);

CREATE TABLE IF NOT EXISTS gifts (
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  id text PRIMARY KEY,
  gift_expectation_id text NOT NULL,
  amount integer
);

CREATE TABLE IF NOT EXISTS receivers (
  id text PRIMARY KEY,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  gift text,
  first_name text,
  last_name text
);

CREATE TABLE IF NOT EXISTS questionnaire_sof_giftors (
  id text PRIMARY KEY
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS questionnaire_sofs (
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  id text PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS properties (
  id text PRIMARY KEY,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  property_expectation_id text NOT NULL,
  address text,
  new_build boolean,
  price integer,
  stamp_duty integer
);

CREATE TABLE IF NOT EXISTS identity_documents (
  id text PRIMARY KEY,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  document_id text,
  type text,
  expiry_date text,
  issuing_country text,
  mrz text[],
  nfc_chip_enabled boolean,
  readid_id text
);


CREATE TABLE IF NOT EXISTS readids (
  id text PRIMARY KEY,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  session text,
  nfc_expectation_id text,
  documents_identity_expectation_id text
);

CREATE TABLE IF NOT EXISTS iproovs (
  id text PRIMARY KEY,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  enrolment_token_created_at timestamptz,
  enrolment_token_value text,
  token_created_at timestamptz,
  token_value text,
  iproov_biometrics_expectation_id text
);

CREATE TABLE IF NOT EXISTS iproov_attempts (
  id text PRIMARY KEY,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  frame text,
  frame_available boolean,
  passed boolean,
  status text,
  token_created_at timestamptz,
  token_validated_at timestamptz,
  token_value text,
  iproov_id text,
  iproov_biometrics_expectation_id text
);