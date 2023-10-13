CREATE TABLE "product"
(
    "id"          serial PRIMARY KEY,
    "updated_at"  timestamptz NOT NULL DEFAULT (now()),
    "created_at"  timestamptz NOT NULL DEFAULT (now()),
    "name"        text NOT NULL UNIQUE,
    "slug"        text NOT NULL UNIQUE,
    "description" text NOT NULL DEFAULT (''),
    "price"       int NOT NULL ,
    "images"      text[],


    "category_id" int NOT NULL
);

ALTER TABLE "product"
    ADD FOREIGN KEY ("category_id") REFERENCES "category" ("id");