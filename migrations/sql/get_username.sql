CREATE OR REPLACE FUNCTION get_username(code INT) 
RETURNS JSON
LANGUAGE plpgsql
AS
$$
DECLARE
    result_record RECORD;
BEGIN
   SELECT cr.record_code, u.user_name
    INTO result_record
    FROM users u
    JOIN credentials cr ON u.credentials = cr.id
    WHERE cr.record_code = Code;

    RETURN JSON_BUILD_OBJECT(
        'record_code', result_record.record_code,
        'user_name', result_record.user_name
    );
END;
$$;
