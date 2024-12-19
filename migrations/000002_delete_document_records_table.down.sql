-- Down migration for dropping the "document_records" table and foreign key constraint

-- Drop the foreign key constraint
ALTER TABLE "document_records" DROP CONSTRAINT "document_records_upload_by_id_fkey";

-- Drop the "document_records" table
DROP TABLE IF EXISTS "document_records";