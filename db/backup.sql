-- Create database if it doesn't already exist
DO $$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'trullio') THEN
      CREATE DATABASE trullio;
   END IF;
END $$;

-- Create user if it doesn't already exist
DO $$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'caciano4') THEN
      CREATE USER caciano4 WITH ENCRYPTED PASSWORD 'caciano123';
   END IF;
END $$;

-- Grant privileges if not already granted
GRANT ALL PRIVILEGES ON DATABASE trullio TO caciano;
