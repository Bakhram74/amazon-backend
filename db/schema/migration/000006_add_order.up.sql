CREATE TYPE enum_order_status AS ENUM ('PENDING', 'PAYED', 'SHIPPED','DELIVERED');

CREATE TABLE "order"
(
    "id"         serial PRIMARY KEY,
    "updated_at" timestamptz,
    "created_at" timestamptz       DEFAULT (now()),
    "status"     enum_order_status DEFAULT ('PENDING'),

    "user_id"    bigint
);



ALTER TABLE "order"
    ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");
