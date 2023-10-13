
CREATE TABLE "order"
(
    "id"         serial PRIMARY KEY,
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "created_at" timestamptz   NOT NULL    DEFAULT (now()),
    "status"     enum_order_status NOT NULL DEFAULT ('PENDING'),

    "user_id"    bigint NOT NULL
);



ALTER TABLE "order"
    ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");
