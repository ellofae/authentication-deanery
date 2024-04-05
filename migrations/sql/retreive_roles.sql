CREATE OR REPLACE FUNCTION retreive_roles() 
RETURNS JSON
LANGUAGE plpgsql
AS
$$
DECLARE
    result JSON;
BEGIN
    SELECT json_agg(row_to_json(existing_roles)) INTO result
    FROM (
        SELECT status_name FROM public."user_status"
    ) AS existing_roles;
    
    RETURN result;
END;
$$;
