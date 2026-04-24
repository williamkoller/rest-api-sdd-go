CREATE TYPE class_shift AS ENUM ('morning', 'afternoon', 'full');

CREATE TABLE classes (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    unit_id       UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    classroom_id  UUID REFERENCES classrooms(id) ON DELETE SET NULL,
    name          VARCHAR(100) NOT NULL,
    grade_level   VARCHAR(20) NOT NULL,
    shift         class_shift NOT NULL,
    academic_year INTEGER NOT NULL,
    active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (unit_id, name, academic_year)
);

CREATE INDEX idx_classes_unit_id ON classes(unit_id);
