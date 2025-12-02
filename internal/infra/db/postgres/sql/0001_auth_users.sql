CREATE TYPE auth_role AS ENUM ('admin', 'user');

CREATE TABLE IF NOT EXISTS auth_users
(
    id           UUID PRIMARY KEY             DEFAULT gen_random_uuid(),
    created_at   TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,

    username     VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255)        NOT NULL,
    role         auth_role           NOT NULL
);

CREATE RULE auth_users_disable_delete AS ON DELETE TO auth_users DO INSTEAD NOTHING;