CREATE TABLE "product"
(
    "id"          serial PRIMARY KEY,
    "updated_at"  timestamptz,
    "created_at"  timestamptz DEFAULT (now()),
    "name"        text UNIQUE,
    "slug"        text UNIQUE,
    "description" text,
    "price"       int,
    "images"      text[],


    "category_id" int
);

ALTER TABLE "product"
    ADD FOREIGN KEY ("category_id") REFERENCES "category" ("id");