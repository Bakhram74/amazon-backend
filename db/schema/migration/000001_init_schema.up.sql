CREATE TABLE "user"
(
    "id"          bigserial PRIMARY KEY,
    "name"        text        NOT NULL UNIQUE,
    "email"       text        NOT NULL UNIQUE,
    "phone"       text        NOT NULL DEFAULT (''),
    "password"    text        NOT NULL,
    "avatar_path" text        NOT NULL DEFAULT ('/uploads/default-avatar.png'),
    "updated_at"  timestamptz ,
    "created_at"  timestamptz NOT NULL DEFAULT (now())
);
