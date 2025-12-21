DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'ddfilms_user') THEN
        CREATE USER ddfilms_user;
        RAISE NOTICE 'User ddfilms_user created';
    ELSE
        RAISE NOTICE 'User ddfilms_user already exists';
    END IF;
END
$$;