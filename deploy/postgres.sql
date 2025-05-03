------------------------------------------------------------
CREATE EXTENSION IF NOT EXISTS btree_gist;
------------------------------------------------------------
-- Дома
CREATE TABLE houses (
   id smallserial PRIMARY KEY,
   name text NOT NULL,
   slug text NOT NULL UNIQUE,
   capacity smallint NOT NULL CHECK (capacity > 0),
   base_price numeric(10,2) NOT NULL CHECK (base_price >= 0),
   description text,
   created_at timestamptz NOT NULL DEFAULT now(),
   updated_at timestamptz NOT NULL DEFAULT now()
);
------------------------------------------------------------
-- Гости
CREATE TABLE guests (
   id bigserial PRIMARY KEY,
   first_name text NOT NULL,
   last_name text NOT NULL,
   email citext NOT NULL UNIQUE,
   phone text,
   created_at timestamptz NOT NULL DEFAULT now()
);
------------------------------------------------------------
-- Статусы
CREATE TYPE reservation_status AS ENUM (
    'pending','confirmed','checked_in','checked_out','cancelled'
);

CREATE TYPE payment_status AS ENUM (
    'pending','paid','failed','refunded'
);
------------------------------------------------------------
-- Брони
CREATE TABLE reservations (
    id bigserial PRIMARY KEY,
    house_id smallint REFERENCES houses ON DELETE CASCADE,
    guest_id bigint REFERENCES guests ON DELETE RESTRICT,
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
CREATE TABLE payments (
    id bigserial PRIMARY KEY,
    reservation_id bigint REFERENCES reservations ON DELETE CASCADE,
    amount numeric(10,2) NOT NULL,
    currency char(3) DEFAULT 'RUB',
    method text,
    status payment_status NOT NULL,
    gateway_tx_id text, -- id в платёжном шлюзе
    paid_at timestamptz
);
------------------------------------------------------------
-- Доп. услуги
CREATE TABLE extras (
    id serial PRIMARY KEY,
    code text UNIQUE,
    name text NOT NULL,
    price numeric(10,2) NOT NULL
);

CREATE TABLE reservation_extras (
    reservation_id bigint REFERENCES reservations ON DELETE CASCADE,
    extra_id int  REFERENCES extras ON DELETE RESTRICT,
    quantity smallint NOT NULL DEFAULT 1,
    amount numeric(10,2) NOT NULL,
    PRIMARY KEY (reservation_id, extra_id)
);
------------------------------------------------------------
-- Блокировка дат (ремонт, частное пользование)
CREATE TABLE blackouts (
    id serial PRIMARY KEY,
    house_id smallint REFERENCES houses ON DELETE CASCADE,
    period daterange NOT NULL,
    reason text,
    CONSTRAINT no_overlap_blackout
    EXCLUDE USING gist (house_id WITH =, period WITH &&)
);
------------------------------------------------------------
CREATE INDEX reservations_active_idx
    ON reservations
    USING gist (house_id, stay);
------------------------------------------------------------