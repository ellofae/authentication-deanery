CREATE OR REPLACE FUNCTION set_user_password(credentials_id INT, encrypted_password TEXT) 
RETURNS VOID
LANGUAGE plpgsql
AS
$$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables 
    WHERE table_schema = 'public' AND table_name = 'credentials') THEN
        UPDATE credentials SET user_password = encrypted_password
        WHERE id = credentials_id;
    ELSE
        RAISE EXCEPTION 'Relation credentials is not present.';
    END IF;
END;
$$;