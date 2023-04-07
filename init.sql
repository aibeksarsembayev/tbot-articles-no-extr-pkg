CREATE TABLE IF NOT EXISTS "article"(  
  "article_id" BIGINT NOT NULL UNIQUE PRIMARY KEY,
  "title" TEXT NOT NULL,     
  "created_at" TIMESTAMP NOT NULL,
  "updated_at" TIMESTAMP NOT NULL,  
  "url" TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "user" (
  "user_id" BIGINT PRIMARY KEY,
  "author_first_name" TEXT NOT NULL,
  "author_last_name" TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "category"(
  "category_id" BIGSERIAL PRIMARY KEY,
  "category_name" TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "article_user"(
  "article_id" BIGINT NOT NULL REFERENCES "article" ("article_id") ON DELETE CASCADE,
  "user_id" BIGINT NOT NULL REFERENCES "user" ("user_id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "article_category"(
  "article_id" BIGINT NOT NULL REFERENCES "article" ("article_id") ON DELETE CASCADE,
  "category_id" BIGINT NOT NULL REFERENCES "category" ("category_id") ON DELETE CASCADE
);