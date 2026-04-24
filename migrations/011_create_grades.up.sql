CREATE TABLE grades (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    enrollment_id UUID NOT NULL REFERENCES enrollments(id) ON DELETE CASCADE,
    subject       VARCHAR(100) NOT NULL,
    period        VARCHAR(20) NOT NULL,
    value         DECIMAL(5,2) NOT NULL CHECK (value >= 0 AND value <= 10),
    recorded_by   UUID NOT NULL REFERENCES users(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (enrollment_id, subject, period)
);

CREATE INDEX idx_grades_enrollment_id ON grades(enrollment_id);
