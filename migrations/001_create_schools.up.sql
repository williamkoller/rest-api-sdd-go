CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE schools (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(200) NOT NULL,
    cnpj       VARCHAR(14)  NOT NULL UNIQUE,
    email      VARCHAR(255) NOT NULL,
    phone      VARCHAR(20),
    active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
