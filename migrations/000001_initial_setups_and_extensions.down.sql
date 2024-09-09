-- Set time zone of database to Ashgabat
ALTER DATABASE ecommerce SET TIME ZONE 'Asia/Ashgabat';

-- Create uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create case-insensitive character string type
CREATE EXTENSION IF NOT EXISTS citext;