-- Revert time zone to the default PostgreSQL time zone (typically UTC)
ALTER DATABASE ecommerce SET TIME ZONE 'UTC';