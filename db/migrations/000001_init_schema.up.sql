CREATE TABLE "keypool" (
  "id" bigserial PRIMARY KEY,
  "shortUrl" varchar,
  "created_at" timestamp NOT NULL,
  "is_used" bool
);

CREATE TABLE "keylink" (
  "id" bigserial PRIMARY KEY,
  "shortUrl" varchar,
  "originalUrl" text,
  "created_at" timestamp NOT NULL
);

CREATE INDEX ON "keypool" ("is_used");

ALTER TABLE "keypool" ADD FOREIGN KEY ("shortUrl") REFERENCES "keylink" ("shortUrl") DEFERRABLE INITIALLY IMMEDIATE;