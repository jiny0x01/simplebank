CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);


ALTER TABLE "account" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");


-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
-- 각 owner는 1개의 account당 하나의 currency만 갖게하는 쿼리
-- 위와 동일하게 만드는 쿼리로 unique constraint를 account의 owner와 currency에 부여한다.

ALTER TABLE "account" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");