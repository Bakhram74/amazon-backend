CREATE TABLE "order_item"
(
    "id"         serial PRIMARY KEY,
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "quantity"   int NOT NULL ,
    "price"      int NOT NULL ,

    "order_id"   int NOT NULL ,
    "product_id" int NOT NULL
);

ALTER TABLE "order_item"
    ADD FOREIGN KEY ("order_id") REFERENCES "order" ("id");
ALTER TABLE "order_item"
    ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");