CREATE OR REPLACE FUNCTION create_user(display_name text, email_arg text, phone_number text, status_arg text)
RETURNS JSON
LANGUAGE plpgsql
AS
$$
DECLARE
    result_json JSON;
    credentials_id INT;
    status_name_val TEXT;
    record_code_registered INT;
    current_datetime TIMESTAMP;
    phone_exists BOOLEAN;
    email_exists BOOLEAN;
    record_exists BOOLEAN;
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables 
    WHERE table_schema = 'public' AND table_name = 'credentials') THEN

        SELECT EXISTS (SELECT 1 FROM "credentials" WHERE phone = phone_number) INTO phone_exists FOR UPDATE;
        SELECT EXISTS (SELECT 1 FROM "credentials" WHERE email = email_arg) INTO email_exists FOR UPDATE;

        IF phone_exists OR email_exists THEN
            RAISE EXCEPTION 'Phone or Email is already in use.';
        END IF;

        IF EXISTS (SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'student_record') THEN

            SELECT EXISTS (SELECT 1 FROM student_record WHERE phone = phone_number) INTO record_exists FOR UPDATE;

            IF record_exists THEN
                RAISE EXCEPTION 'Record code for this number is already registered.';
            END IF;

            INSERT INTO student_record(phone) VALUES (phone_number) RETURNING record_code INTO record_code_registered;
        ELSE
            RAISE EXCEPTION 'Relation student_record is not present.';
        END IF;

        SELECT current_timestamp INTO current_datetime;
        INSERT INTO "credentials"(email, phone, register_date, record_code) VALUES (email_arg, phone_number, current_datetime, record_code_registered) RETURNING id INTO credentials_id;
    ELSE
        RAISE EXCEPTION 'Relation credentials is not present.';
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.tables 
    WHERE table_schema = 'public' AND table_name = 'users') THEN
        INSERT INTO "users"(user_name, user_status, credentials) VALUES (display_name, status_arg, credentials_id);
    ELSE
        RAISE EXCEPTION 'Relation user is not present.';
    END IF;

    result_json = json_build_object('user_name', display_name, 'phone', phone_number, 'email', email_arg, 'credentials', credentials_id, 'user_status', status_arg, 'register_date', current_datetime::text);
    RETURN result_json;
END;
$$;