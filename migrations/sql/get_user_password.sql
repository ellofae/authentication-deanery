CREATE OR REPLACE FUNCTION get_user_password(code INT) 
RETURNS TEXT
LANGUAGE plpgsql
AS
$$
DECLARE
    received_password TEXT;
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables 
    WHERE table_schema = 'public' AND table_name = 'credentials') THEN
        SELECT user_password INTO received_password FROM credentials WHERE record_code = code;
    ELSE
        RAISE EXCEPTION 'Relation credentials is not present.';
    END IF;

    RETURN received_password;
END;
$$;
