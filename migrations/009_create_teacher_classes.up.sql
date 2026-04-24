CREATE TABLE teacher_classes (
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    class_id   UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject    VARCHAR(100) NOT NULL,
    PRIMARY KEY (teacher_id, class_id, subject)
);
