-- migrate:up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TABLE account (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    first_name varchar(250) NOT NULL,
    last_name varchar(250) NOT NULL,
    email varchar(250) UNIQUE NOT NULL,
    dob date NOT NULL, 
    gender gender NOT NULL,
    phone_number varchar(250) NOT NULL,
    password text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS account;
DROP TYPE IF EXISTS gender;
