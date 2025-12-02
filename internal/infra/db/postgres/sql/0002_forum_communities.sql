CREATE TABLE forum_communities
(
    id          UUID PRIMARY KEY             DEFAULT gen_random_uuid(),
    created_at  TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,

    name        VARCHAR(255) UNIQUE NOT NULL,
    description VARCHAR(255)        NOT NULL
);

CREATE RULE forum_communities_disable_delete AS ON DELETE TO forum_communities DO INSTEAD NOTHING;