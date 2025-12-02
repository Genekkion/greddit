CREATE TABLE forum_posts
(
    id           UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,

    title        VARCHAR(255) NOT NULL,
    body         TEXT         NOT NULL,

    poster_id    UUID         NOT NULL REFERENCES auth_users (id) ON DELETE CASCADE,
    community_id UUID         NOT NULL REFERENCES forum_communities (id) ON DELETE CASCADE
);

CREATE RULE forum_posts_disable_delete AS ON DELETE TO forum_posts DO INSTEAD NOTHING;