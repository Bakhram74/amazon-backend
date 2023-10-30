CREATE TABLE "category"
(
    "id"        INTEGER PRIMARY KEY,
    "updated_at"  timestamptz NOT NULL DEFAULT (now()),
    "created_at"  timestamptz NOT NULL  DEFAULT (now()),
    "name"        text NOT NULL UNIQUE
);