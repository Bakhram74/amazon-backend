CREATE TABLE "sessions"
(
    "id" uuid PRIMARY KEY,
    "refresh_token" varchar     NOT NULL,
    "user_agent"    varchar     NOT NULL,
    "client_ip"     varchar     NOT NULL,
    "is_blocked"    boolean     NOT NULL DEFAULT false,
    "expires_at"    timestamptz NOT NULL,
    "created_at"    timestamptz NOT NULL DEFAULT (now()),

    "user_id"          bigint     NOT NULL
);

ALTER TABLE "sessions"
    ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");