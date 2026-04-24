CREATE TYPE waitlist_status AS ENUM ('waiting', 'offer_made', 'accepted', 'declined', 'expired');

CREATE TABLE waitlist_entries (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    unit_id       UUID NOT NULL REFERENCES units(id),
    guardian_name  VARCHAR(200) NOT NULL,
    guardian_email VARCHAR(200) NOT NULL,
    student_name  VARCHAR(200) NOT NULL,
    grade_level   VARCHAR(50) NOT NULL,
    academic_year INTEGER NOT NULL,
    position      INTEGER NOT NULL,
    status        waitlist_status NOT NULL DEFAULT 'waiting',
    referral_id   UUID REFERENCES waitlist_entries(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_waitlist_unit_status ON waitlist_entries(unit_id, status);
CREATE INDEX idx_waitlist_unit_grade_year ON waitlist_entries(unit_id, grade_level, academic_year);
