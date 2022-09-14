CREATE TABLE "oauths" (
  "id" varchar PRIMARY KEY,
  "fullname" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "provider" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);


