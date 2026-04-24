CREATE TABLE guardian_students (
    guardian_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    student_id   UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    relationship VARCHAR(50),
    PRIMARY KEY (guardian_id, student_id)
);
