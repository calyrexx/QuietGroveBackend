------------------------------------------------------------
CREATE EXTENSION IF NOT EXISTS btree_gist;
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
    email text NOT NULL,
    phone text,
    tg_user_id   bigint,
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
-- Баня\чан
CREATE TABLE IF NOT EXISTS sauna_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    default_price NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS house_saunas (
    id SERIAL PRIMARY KEY,
    house_id SMALLINT NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
    sauna_type_id INT NOT NULL REFERENCES sauna_types(id) ON DELETE CASCADE,
    price NUMERIC(10,2) NOT NULL,
    description TEXT,
    images TEXT[] NOT NULL DEFAULT '{}',
    is_active BOOL NOT NULL DEFAULT TRUE,
    UNIQUE (house_id, sauna_type_id)
);

CREATE TABLE IF NOT EXISTS sauna_fill_options (
    id SERIAL PRIMARY KEY,
    sauna_type_id INT NOT NULL REFERENCES sauna_types(id) ON DELETE CASCADE,
    label TEXT NOT NULL,
    image TEXT NOT NULL,
    description TEXT NOT NULL,
    price NUMERIC(10,2),
    is_active BOOL NOT NULL DEFAULT TRUE
);
------------------------------------------------------------
-- Брони бани\чана
CREATE TABLE IF NOT EXISTS bathhouse_reservations (
    id serial PRIMARY KEY,
    reservation_uuid uuid NOT NULL REFERENCES reservations(uuid) ON DELETE CASCADE,
    type_id int NOT NULL REFERENCES sauna_types(id) ON DELETE CASCADE, -- Баня или чан
    date date NOT NULL,
    time_from time NOT NULL,
    time_to time NOT NULL,
    fill_option_id int REFERENCES sauna_fill_options(id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT bathhouse_time_unique UNIQUE (type_id, date, time_from, time_to)
);
------------------------------------------------------------
-- Доп. услуги
CREATE TABLE IF NOT EXISTS extras (
    id smallint PRIMARY KEY,
    name text NOT NULL,
    short_text text NOT NULL,
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
-- Верификация пользователя
CREATE TABLE IF NOT EXISTS verifications (
    uuid         uuid  PRIMARY KEY,
    code         char(6)    NOT NULL,
    email        text       NOT NULL,
    phone        text       NOT NULL,
    name         text       NOT NULL,
    tg_user_id   bigint,
    status       text       NOT NULL,
    created_at   timestamptz DEFAULT now(),
    verified_at  timestamptz,
    expires_at   timestamptz NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_verifications_code
    ON verifications(code);
------------------------------------------------------------
CREATE INDEX IF NOT EXISTS reservations_active_idx
    ON reservations
    USING gist (house_id, stay);
------------------------------------------------------------