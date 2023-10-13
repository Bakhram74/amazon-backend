CREATE TABLE "review"
(
    "id"          serial PRIMARY KEY,
    "updated_at"  timestamptz NOT NULL DEFAULT (now()),
    "created_at"  timestamptz NOT NULL  DEFAULT (now()),
    "rating" int NOT NULL DEFAULT (0),
    "text" text NOT NULL DEFAULT (''),

    "user_id" bigint NOT NULL ,
    "product_id" int NOT NULL
);

ALTER TABLE "review"
    ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");
ALTER TABLE "review"
    ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");