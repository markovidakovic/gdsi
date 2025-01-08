-- migrate:up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TABLE account (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name varchar(50) NOT NULL,
    last_name varchar(50) NOT NULL,
    email varchar(100) NOT NULL UNIQUE,
    dob date NOT NULL, 
    gender gender NOT NULL,
    phone_number varchar(50) NOT NULL,
    password text NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS account;
DROP TYPE IF EXISTS gender;
