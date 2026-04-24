CREATE TABLE feed_posts (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id    UUID NOT NULL REFERENCES schools(id),
    unit_id      UUID REFERENCES units(id),
    author_id    UUID NOT NULL REFERENCES users(id),
    body         TEXT NOT NULL,
    image_url    VARCHAR(500),
    published_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_feed_posts_school_published ON feed_posts(school_id, published_at DESC);
