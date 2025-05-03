------------------------------------------------------------

CREATE TABLE IF NOT EXISTS reservations (
    uuid uuid NOT NULL,
    added_at BIGINT NOT NULL,
    CONSTRAINT reservations_pk PRIMARY KEY (uuid),
);

------------------------------------------------------------