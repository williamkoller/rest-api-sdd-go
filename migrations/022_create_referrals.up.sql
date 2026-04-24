CREATE TYPE referral_status AS ENUM ('pending', 'registered', 'enrolled');

CREATE TABLE referrals (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id      UUID NOT NULL REFERENCES schools(id),
    referrer_id    UUID NOT NULL REFERENCES users(id),
    referral_code  VARCHAR(8) NOT NULL UNIQUE,
    referred_email VARCHAR(200),
    status         referral_status NOT NULL DEFAULT 'pending',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_referrals_school_referrer ON referrals(school_id, referrer_id);
