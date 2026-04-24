CREATE TABLE curriculum_entries (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    class_id   UUID NOT NULL REFERENCES classes(id),
    subject    VARCHAR(100) NOT NULL,
    teacher_id UUID NOT NULL REFERENCES users(id),
    day_of_week day_of_week NOT NULL,
    start_time TIME NOT NULL,
    end_time   TIME NOT NULL,
    UNIQUE (class_id, day_of_week, start_time)
);

CREATE INDEX idx_curriculum_class_id ON curriculum_entries(class_id);
