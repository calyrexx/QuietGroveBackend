------------------------------------------------------------
CREATE EXTENSION IF NOT EXISTS btree_gist;
CREATE EXTENSION IF NOT EXISTS citext;
------------------------------------------------------------
-- Дома
CREATE TABLE IF NOT EXISTS houses (
    id smallint PRIMARY KEY,
    name text NOT NULL,
    capacity smallint NOT NULL CHECK (capacity > 0),
    base_price numeric(10,2) NOT NULL CHECK (base_price >= 0),
    description text,
    images text[] NOT NULL DEFAULT '{}'::text[],
    check_in_from text NOT NULL DEFAULT '14:00',
    check_out_until text NOT NULL DEFAULT '11:00',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
------------------------------------------------------------
-- Гости
CREATE TABLE IF NOT EXISTS guests (
    uuid uuid PRIMARY KEY,
    name text NOT NULL,
    email citext UNIQUE,
    phone text,
    created_at timestamptz NOT NULL DEFAULT now()
);
------------------------------------------------------------
-- Статусы
DO $$
    BEGIN
        CREATE TYPE reservation_status AS ENUM
            ('pending','confirmed','checked_in','checked_out','cancelled');
    EXCEPTION
        WHEN duplicate_object THEN NULL;
END $$;

DO $$
    BEGIN
        CREATE TYPE payment_status AS ENUM
            ('pending','paid','failed','refunded');
    EXCEPTION
        WHEN duplicate_object THEN NULL;
END $$;
------------------------------------------------------------
-- Брони
CREATE TABLE IF NOT EXISTS reservations (
    uuid uuid PRIMARY KEY,
    house_id smallint REFERENCES houses ON DELETE CASCADE,
    guest_uuid uuid REFERENCES guests ON DELETE RESTRICT,
    stay daterange NOT NULL, -- диапазон дат заезда ([))
    guests_count smallint  NOT NULL CHECK (guests_count > 0),
    status reservation_status NOT NULL DEFAULT 'pending',
    total_price numeric(10,2) NOT NULL CHECK (total_price >= 0),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT no_overlap
    EXCLUDE USING gist (
        house_id WITH =,
        stay WITH &&  -- «&&» — пересечение диапазонов
    )
);
------------------------------------------------------------
-- Платежи
CREATE TABLE IF NOT EXISTS payments (
    uuid uuid PRIMARY KEY,
    reservation_uuid uuid REFERENCES reservations ON DELETE CASCADE,
    amount numeric(10,2) NOT NULL,
    currency char(3) DEFAULT 'RUB',
    method text,
    status payment_status NOT NULL,
    gateway_tx_id text, -- id в платёжном шлюзе
    paid_at timestamptz
);
------------------------------------------------------------
-- Доп. услуги
CREATE TABLE IF NOT EXISTS extras (
    id smallint PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL,
    images text[] NOT NULL DEFAULT '{}'::text[],
    price numeric(10,2) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS reservation_extras (
    reservation_uuid uuid REFERENCES reservations ON DELETE CASCADE,
    extra_id int REFERENCES extras ON DELETE RESTRICT,
    quantity smallint NOT NULL DEFAULT 1,
    amount numeric(10,2) NOT NULL,
    PRIMARY KEY (reservation_uuid, extra_id)
);
------------------------------------------------------------
-- Блокировка дат (ремонт, частное пользование)
CREATE TABLE IF NOT EXISTS blackouts (
    id serial PRIMARY KEY,
    house_id smallint REFERENCES houses ON DELETE CASCADE,
    period daterange NOT NULL,
    reason text,
    CONSTRAINT no_overlap_blackout
    EXCLUDE USING gist (
        house_id WITH =,
        period WITH &&
    )
);
------------------------------------------------------------
CREATE INDEX IF NOT EXISTS reservations_active_idx
    ON reservations
    USING gist (house_id, stay);
------------------------------------------------------------