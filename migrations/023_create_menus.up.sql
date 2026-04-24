CREATE TYPE day_of_week AS ENUM ('monday', 'tuesday', 'wednesday', 'thursday', 'friday');
CREATE TYPE meal_type AS ENUM ('breakfast', 'lunch', 'snack', 'dinner');

CREATE TABLE menus (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    unit_id    UUID NOT NULL REFERENCES units(id),
    week_start DATE NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (unit_id, week_start)
);

CREATE TABLE menu_items (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    menu_id     UUID NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
    day_of_week day_of_week NOT NULL,
    meal_type   meal_type NOT NULL,
    description TEXT NOT NULL
);

CREATE INDEX idx_menu_items_menu_id ON menu_items(menu_id);
