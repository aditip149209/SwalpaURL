CREATE TABLE "key_pool" (
  "id" bigserial PRIMARY KEY,
  "short_url" varchar NOT NULL UNIQUE,
  "created_at" timestamp NOT NULL DEFAULT NOW(),
  "is_used" bool DEFAULT false
);

CREATE TABLE "key_link" (
  "id" bigserial PRIMARY KEY,
  "short_url" varchar UNIQUE NOT NULL,
  "original_url" text NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT NOW(),
  FOREIGN KEY ("short_url") REFERENCES "key_pool" ("short_url") ON DELETE CASCADE
);

CREATE INDEX ON "key_pool" ("is_used");
CREATE INDEX ON "key_pool" ("short_url");