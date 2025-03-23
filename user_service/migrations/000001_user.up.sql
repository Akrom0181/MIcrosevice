CREATE TYPE user_type AS ENUM (
  'user',
  'admin'
);

CREATE TYPE user_role AS ENUM (
  'user',
  'admin'
);

CREATE TYPE gender AS ENUM (
  'male',
  'female'
);

CREATE TYPE user_status AS ENUM (
  'active',
  'blocked',
  'inverify'
);

CREATE TYPE platform AS ENUM (
  'admin',
  'web',
  'mobile'
);

CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY,
  user_type user_type NOT NULL,
  user_role user_role NOT NULL,
  full_name VARCHAR(100) NOT NULL,
  user_name VARCHAR(50) UNIQUE NOT NULL,
  email VARCHAR(50) UNIQUE NOT NULL,
  password VARCHAR(150) NOT NULL,
  gender gender NOT NULL DEFAULT 'male',
  status user_status NOT NULL DEFAULT 'inverify',
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);


CREATE TABLE IF NOT EXISTS session (
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  user_agent text NOT NULL,
  is_active boolean NOT NULL,
  expires_at timestamp,
  last_active_at timestamp,
  platform platform NOT NULL,
  ip_address varchar(64) NOT NULL,
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);
