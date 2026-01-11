BEGIN;
CREATE TABLE IF NOT EXISTS decimal_exp
(
    id              UUID      NOT NULL,
    balance_a       NUMERIC   NOT NULL,
    balance_b       NUMERIC   NULL,
    balance_history NUMERIC[] NOT NULL,
    CONSTRAINT decimal_exp_pk PRIMARY KEY (id)
);
COMMIT;