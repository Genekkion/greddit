CREATE TABLE forum_comments
(
    id           UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,

    body         TEXT        NOT NULL,

    commenter_id UUID        NOT NULL REFERENCES auth_users (id) ON DELETE CASCADE,
    post_id      UUID        NOT NULL REFERENCES forum_posts (id) ON DELETE CASCADE,
    parent_id    UUID REFERENCES forum_comments (id)
);

CREATE RULE forum_comments_disable_delete AS ON DELETE TO forum_comments DO INSTEAD NOTHING;