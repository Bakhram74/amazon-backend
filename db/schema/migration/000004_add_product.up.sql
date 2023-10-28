CREATE TABLE "product"
(
    "id"          serial PRIMARY KEY,
    "updated_at"  timestamptz NOT NULL DEFAULT (now()),
    "created_at"  timestamptz NOT NULL DEFAULT (now()),
    "name"        text        NOT NULL UNIQUE,
    "slug"        text        NOT NULL UNIQUE,
    "description" text        NOT NULL DEFAULT (''),
    "price"       int         NOT NULL,
    "images"      text[],

    "user_id"     bigint      NOT NULL DEFAULT (0),
    "category_id" int         NOT NULL DEFAULT (0)
);

ALTER TABLE "product"
    ADD FOREIGN KEY ("category_id") REFERENCES "category" ("id");
ALTER TABLE "product"
    ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");