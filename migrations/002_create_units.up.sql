CREATE TABLE units (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id  UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name       VARCHAR(200) NOT NULL,
    address    TEXT NOT NULL,
    city       VARCHAR(100) NOT NULL,
    state      CHAR(2) NOT NULL,
    zip_code   VARCHAR(10) NOT NULL,
    phone      VARCHAR(20),
    active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (school_id, name)
);
