CREATE TYPE attendance_status AS ENUM ('present', 'absent', 'justified');

CREATE TABLE attendance_records (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    enrollment_id UUID NOT NULL REFERENCES enrollments(id) ON DELETE CASCADE,
    date          DATE NOT NULL,
    status        attendance_status NOT NULL,
    note          TEXT,
    recorded_by   UUID NOT NULL REFERENCES users(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (enrollment_id, date)
);

CREATE INDEX idx_attendance_enrollment_id ON attendance_records(enrollment_id);
