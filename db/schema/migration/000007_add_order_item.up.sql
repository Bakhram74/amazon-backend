CREATE TABLE "order_item"
(
    "id"         serial PRIMARY KEY,
    "updated_at" timestamptz,
    "created_at" timestamptz DEFAULT (now()),
    "quantity"   int,
    "price"      int,

    "order_id"   int,
    "product_id" int
);

ALTER TABLE "order_item"
    ADD FOREIGN KEY ("order_id") REFERENCES "order" ("id");
ALTER TABLE "order_item"
    ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");