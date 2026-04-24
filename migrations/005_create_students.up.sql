CREATE TABLE students (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id           UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name                VARCHAR(200) NOT NULL,
    birth_date          DATE NOT NULL,
    cpf                 VARCHAR(11),
    registration_number VARCHAR(50) NOT NULL,
    active              BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (school_id, registration_number)
);

CREATE INDEX idx_students_school_id ON students(school_id);
