CREATE TYPE invoice_status AS ENUM ('pending', 'paid', 'overdue', 'cancelled');
CREATE TYPE payment_method AS ENUM ('pix', 'boleto', 'credit_card', 'manual');

CREATE TABLE invoices (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    school_id  UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    amount     DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    due_date   DATE NOT NULL,
    reference  VARCHAR(20) NOT NULL,
    status     invoice_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (student_id, reference)
);

CREATE TABLE payments (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invoice_id  UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    amount_paid DECIMAL(10,2) NOT NULL,
    paid_at     TIMESTAMPTZ NOT NULL,
    method      payment_method NOT NULL,
    gateway_ref VARCHAR(255),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_invoices_student_id ON invoices(student_id);
CREATE INDEX idx_invoices_school_id ON invoices(school_id);
CREATE INDEX idx_invoices_status ON invoices(status);
