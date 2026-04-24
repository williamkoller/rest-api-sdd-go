CREATE TABLE classrooms (
    id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    unit_id  UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    code     VARCHAR(20) NOT NULL,
    capacity INTEGER NOT NULL CHECK (capacity >= 1),
    active   BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (unit_id, code)
);
