CREATE TABLE "review"
(
    "id"          serial PRIMARY KEY,
    "updated_at"  timestamptz ,
    "created_at"  timestamptz   DEFAULT (now()),
    "rating" int,
    "text" text,

    "user_id" bigint,
    "product_id" int
);

ALTER TABLE "review"
    ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");
ALTER TABLE "review"
    ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");