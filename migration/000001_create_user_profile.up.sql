BEGIN;
CREATE TABLE IF NOT EXISTS user_profile
(
    id            UUID        NOT NULL,
    email         TEXT        NOT NULL,
    first_name    TEXT        NOT NULL,
    last_name     TEXT        NOT NULL,
    date_of_birth DATE        NOT NULL,
    create_time   TIMESTAMPTZ NOT NULL,
    update_time   TIMESTAMPTZ NOT NULL,
    version       INTEGER     NOT NULL,
    CONSTRAINT user_profile_pk PRIMARY KEY (id),
    CONSTRAINT user_profile_email_uk UNIQUE (email)
);
COMMIT;