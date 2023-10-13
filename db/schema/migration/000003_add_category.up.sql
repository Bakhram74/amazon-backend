CREATE TABLE "category"
(
    "id"          serial PRIMARY KEY,
    "updated_at"  timestamptz ,
    "created_at"  timestamptz   DEFAULT (now()),
    "name"        text UNIQUE
);