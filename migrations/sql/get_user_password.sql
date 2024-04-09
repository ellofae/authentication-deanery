CREATE OR REPLACE FUNCTION get_user_password(code INT) 
RETURNS JSON
LANGUAGE plpgsql
AS
$$
DECLARE
    received_password TEXT;
    received_status TEXT;
    credentials_id INT;
    result JSON;
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables 
    WHERE table_schema = 'public' AND table_name = 'credentials') THEN
        SELECT user_password, id INTO received_password, credentials_id FROM credentials WHERE record_code = code;
    ELSE
        RAISE EXCEPTION 'Relation credentials is not present.';
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.tables 
    WHERE table_schema = 'public' AND table_name = 'users') THEN
        SELECT user_status INTO received_status FROM users WHERE credentials = credentials_id;
    ELSE
        RAISE EXCEPTION 'Relation users is not present.';
    END IF;

    result := json_build_object('password', received_password, 'status', received_status);

    RETURN result;
END;
$$;