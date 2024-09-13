-- Grant all privileges on the database to the 'dev' user:
GRANT ALL PRIVILEGES ON DATABASE ecommerce TO dev;

-- Grant the 'dev' user permission to create tables in the public schema:
GRANT CREATE ON SCHEMA public TO dev;

-- Optionally, grant usage on the schema (if not granted already):
GRANT USAGE ON SCHEMA public TO dev;


\dn+ public
