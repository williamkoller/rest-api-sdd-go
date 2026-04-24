CREATE TYPE enrollment_status AS ENUM ('active', 'transferred', 'unenrolled');

CREATE TABLE enrollments (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id    UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    class_id      UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    academic_year INTEGER NOT NULL,
    enrolled_at   TIMESTAMPTZ NOT NULL,
    unenrolled_at TIMESTAMPTZ,
    status        enrollment_status NOT NULL DEFAULT 'active'
);

CREATE INDEX idx_enrollments_student_id ON enrollments(student_id);
CREATE INDEX idx_enrollments_class_id ON enrollments(class_id);

-- Only one active enrollment per student per academic year
CREATE UNIQUE INDEX idx_enrollments_active_student_year
    ON enrollments(student_id, academic_year)
    WHERE status = 'active';
