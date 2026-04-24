CREATE TYPE agenda_item_type AS ENUM ('homework', 'event', 'reminder');

CREATE TABLE agenda_items (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    class_id    UUID NOT NULL REFERENCES classes(id),
    created_by  UUID NOT NULL REFERENCES users(id),
    type        agenda_item_type NOT NULL,
    title       VARCHAR(300) NOT NULL,
    description TEXT,
    due_date    TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_agenda_class_id ON agenda_items(class_id);
CREATE INDEX idx_agenda_due_date ON agenda_items(due_date);
