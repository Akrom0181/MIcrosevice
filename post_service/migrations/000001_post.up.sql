CREATE TYPE attachment_type AS ENUM (
    'photo',
    'video'
);

CREATE TYPE post_status AS ENUM (
    'draft',
    'published'
);

CREATE TABLE IF NOT EXISTS posts (
    id uuid PRIMARY KEY,
    owner_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tags json,
    content text NOT NULL,
    status post_status NOT NULL DEFAULT 'draft',
    created_at timestamp NOT NULL DEFAULT 'now()',
    updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE IF NOT EXISTS post_attachment (
    id uuid PRIMARY KEY,
    post_id uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    filepath varchar NOT NULL,
    content_type attachment_type NOT NULL,
    created_at timestamp NOT NULL DEFAULT 'now()',
    updated_at timestamp NOT NULL DEFAULT 'now()'
);