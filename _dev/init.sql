CREATE USER user_profile_role WITH PASSWORD 'postgrespw';

CREATE DATABASE user_profile_local
    WITH
    OWNER = user_profile_role
    ENCODING = 'UTF8';

GRANT ALL PRIVILEGES ON DATABASE user_profile_local TO user_profile_role;