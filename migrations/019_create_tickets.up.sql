CREATE TYPE ticket_category AS ENUM ('general', 'financial', 'academic', 'administrative');
CREATE TYPE ticket_status AS ENUM ('open', 'in_progress', 'resolved', 'closed');

CREATE TABLE tickets (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id    UUID NOT NULL REFERENCES schools(id),
    unit_id      UUID REFERENCES units(id),
    requester_id UUID NOT NULL REFERENCES users(id),
    category     ticket_category NOT NULL,
    status       ticket_status NOT NULL DEFAULT 'open',
    subject      VARCHAR(500) NOT NULL,
    resolved_at  TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tickets_school_status ON tickets(school_id, status);
CREATE INDEX idx_tickets_requester ON tickets(requester_id);
