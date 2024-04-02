CREATE TABLE IF NOT EXISTS "user"(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT NOT NULL,
    credentials INT NOT NULL
);

CREATE TABLE IF NOT EXISTS "status"(
    status_name TEXT PRIMARY KEY,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "credentials"(
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    record_code INT NOT NULL,
    password TEXT NOT NULL,
    register_date TIMESTAMP NOT NULL
);

ALTER TABLE
    "credentials" ADD CONSTRAINT credentials_email_unique UNIQUE(email);

ALTER TABLE
    "credentials" ADD CONSTRAINT credentials_phone_unique UNIQUE(phone);

ALTER TABLE
    "credentials" ADD CONSTRAINT credentials_record_code_unique UNIQUE(record_code);

ALTER TABLE
    "user" ADD CONSTRAINT user_credentials_foreign FOREIGN KEY (credentials) REFERENCES "credentials" (id);

ALTER TABLE
    "user" ADD CONSTRAINT user_status_foreign FOREIGN KEY (status) REFERENCES "status" (status_name);