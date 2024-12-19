-- Down migration for dropping the "users" table and foreign key constraint

-- Drop the "document_records" table
DROP TABLE IF EXISTS "users";