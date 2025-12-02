CREATE USER bigbackend_role WITH PASSWORD 'postgrespw';

CREATE DATABASE bigbackend
    WITH
    OWNER = bigbackend_role
    ENCODING = 'UTF8';

GRANT ALL PRIVILEGES ON DATABASE bigbackend TO bigbackend_role;