CREATE TYPE campaign_status AS ENUM ('open', 'closed');
CREATE TYPE reenrollment_status AS ENUM ('not_started', 'confirmed', 'declined');

CREATE TABLE reenrollment_campaigns (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id     UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    unit_id       UUID REFERENCES units(id) ON DELETE CASCADE,
    academic_year INTEGER NOT NULL,
    deadline      TIMESTAMPTZ NOT NULL,
    status        campaign_status NOT NULL DEFAULT 'open',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE reenrollments (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id  UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    campaign_id UUID NOT NULL REFERENCES reenrollment_campaigns(id) ON DELETE CASCADE,
    status      reenrollment_status NOT NULL DEFAULT 'not_started',
    responded_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (student_id, campaign_id)
);

CREATE INDEX idx_reenrollments_campaign_id ON reenrollments(campaign_id);
