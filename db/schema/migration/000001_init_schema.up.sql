CREATE TABLE "users" (
                         "id" bigserial PRIMARY KEY,
                         "username" text NOT NULL,
                         "phone_number" text NOT NULL,
                         "hashed_password" text NOT NULL,
                         "role" text NOT NULL DEFAULT ('user'),
                         "is_banned" bool NOT NULL DEFAULT (false),
                         "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01',
                         "created_at" timestamptz NOT NULL DEFAULT (now())
);
