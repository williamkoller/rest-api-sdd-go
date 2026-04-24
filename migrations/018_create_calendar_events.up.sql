CREATE TYPE calendar_event_type AS ENUM ('holiday', 'exam_period', 'event', 'recess');

CREATE TABLE calendar_events (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id   UUID NOT NULL REFERENCES schools(id),
    unit_id     UUID REFERENCES units(id),
    title       VARCHAR(300) NOT NULL,
    description TEXT,
    type        calendar_event_type NOT NULL,
    start_date  DATE NOT NULL,
    end_date    DATE NOT NULL,
    created_by  UUID NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT calendar_events_date_order CHECK (end_date >= start_date)
);

CREATE INDEX idx_calendar_school_dates ON calendar_events(school_id, start_date, end_date);
