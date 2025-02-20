CREATE TABLE "document_records" (
  "id" SERIAL PRIMARY KEY,
  "package_file_id" varchar,
  "package_name" varchar,
  "upload_by_id" integer,
  "client_reference_id" varchar UNIQUE,
  "transfer_agent_responsible" varchar,
  "type_of_transfer" varchar,
  "email" varchar,
  "user_id" varchar,
  "first_name" varchar,
  "middle_name" varchar,
  "last_name" varchar,
  "date_of_birth_day" date,
  "personal_phone_number" varchar,
  "street_address" varchar,
  "city" varchar,
  "postal" varchar,
  "letter_state" varchar,
  "letter_country" varchar,
  "national_id" varchar,
  "driver_license" varchar,
  "driver_license_version" varchar,
  "suburb" varchar,
  "voter_id" varchar,
  "passport" varchar, 
  "request" text,
  "response" text,
  "notes" text,
  "match" varchar,
  "complete_kyc" boolean DEFAULT false,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "deleted_at" timestamp
);

ALTER TABLE "document_records" ADD FOREIGN KEY ("upload_by_id") REFERENCES "users" ("id");