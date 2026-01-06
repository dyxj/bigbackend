BEGIN;
CREATE TABLE IF NOT EXISTS user_invitation
(
    id          UUID        NOT NULL,
    email       TEXT        NOT NULL,
    status      TEXT        NOT NULL,
    expiry_time TIMESTAMPTZ NOT NULL,
    token       TEXT        NOT NULL,
    create_time TIMESTAMPTZ NOT NULL,
    update_time TIMESTAMPTZ NOT NULL,
    version     INTEGER     NOT NULL,
    CONSTRAINT user_invitation_pk PRIMARY KEY (id),
    CONSTRAINT user_invitation_token_uk UNIQUE (token)
);
CREATE UNIQUE INDEX user_invitation_accepted_pending_email_uk
    ON user_invitation (email)
    WHERE status IN ('ACCEPTED', 'PENDING');
COMMIT;