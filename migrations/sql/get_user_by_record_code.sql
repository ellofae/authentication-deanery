CREATE OR REPLACE FUNCTION get_credentials_by_record_code(code INT) 
RETURNS JSON
LANGUAGE plpgsql
AS
$$
DECLARE
    result JSON;
BEGIN
    SELECT json_agg(row_to_json(users_with_credentials)) INTO result
    FROM (
        SELECT users.*, credentials.*
        FROM users
        JOIN credentials ON users.credentials = credentials.id
        WHERE credentials.record_code = code
    ) AS users_with_credentials;
    
    RETURN result;
END;
$$;
