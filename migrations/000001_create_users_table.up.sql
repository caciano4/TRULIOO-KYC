CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "phone" varchar,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

INSERT INTO public.users (email, password, first_name, last_name, phone, created_at, updated_at)
VALUES ('caciano4@gmail.com', '123caciano', 'walter', 'caciano', '+55 45 9 99647440', NOW(), NOW());
