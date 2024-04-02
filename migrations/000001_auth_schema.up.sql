CREATE TABLE IF NOT EXISTS "users"(
    id SERIAL PRIMARY KEY,
    user_name TEXT NOT NULL,
    user_status TEXT NOT NULL,
    credentials INT NOT NULL
);

CREATE TABLE IF NOT EXISTS "user_status"(
    status_name TEXT PRIMARY KEY,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "student_record"(
    record_code SERIAL PRIMARY KEY,
    phone TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "credentials"(
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    record_code INT,
    user_password TEXT,
    register_date TIMESTAMP NOT NULL
);

ALTER TABLE
    "credentials" ADD CONSTRAINT credentials_email_unique UNIQUE(email);

ALTER TABLE
    "credentials" ADD CONSTRAINT credentials_phone_unique UNIQUE(phone);

ALTER TABLE
    "student_record" ADD CONSTRAINT student_record_phone_unique UNIQUE(phone);

ALTER TABLE
    "credentials" ADD CONSTRAINT credentials_record_code_unique UNIQUE(record_code);

ALTER TABLE
    "users" ADD CONSTRAINT user_credentials_foreign FOREIGN KEY (credentials) REFERENCES "credentials" (id);

ALTER TABLE
    "users" ADD CONSTRAINT user_status_foreign FOREIGN KEY (user_status) REFERENCES "user_status" (status_name);

ALTER SEQUENCE student_record_record_code_seq RESTART WITH 100000;