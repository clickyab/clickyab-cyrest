
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE SCHEMA aaa;

CREATE SEQUENCE aaa.anon_user START 1001;

CREATE TYPE aaa.user_status AS
ENUM ('registered','verified','banned');

CREATE TABLE aaa.roles(
  id bigserial NOT NULL,
  name varchar NOT NULL,
  description VARCHAR,
  resources VARCHAR[],
  created_at timestamp with time zone DEFAULT NOW(),
  updated_at timestamp with time zone DEFAULT NOW(),
  CONSTRAINT roles_table_primary PRIMARY KEY (id)
);

CREATE TABLE aaa.users(
  id bigserial NOT NULL,
  username varchar NOT NULL,
  password varchar NOT NULL,
  contact varchar NOT NULL,
  token varchar NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  last_login timestamp with time zone DEFAULT NULL,
  status aaa.user_status NOT NULL DEFAULT 'registered',
  attributes jsonb NOT NULL DEFAULT '{}',
  CONSTRAINT users_id_primary PRIMARY KEY (id),
  CONSTRAINT user_name_unique UNIQUE (username),
  CONSTRAINT users_contact_unique UNIQUE (contact)

);

CREATE TABLE aaa.user_roles(
  user_id bigint,
  role_id bigint,
  CONSTRAINT user_role_pk PRIMARY KEY (user_id,role_id)
);

CREATE TABLE aaa.reserved_users(
  id bigserial NOT NULL,
  contact varchar NOT NULL,
  token varchar NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  CONSTRAINT reserved_contacts_primary_id PRIMARY KEY (id),
  CONSTRAINT reserved_code_primary_key UNIQUE (contact)

);

CREATE TABLE aaa.message_logs(
  id bigserial NOT NULL,
  body text NOT NULL,
  contact varchar NOT NULL,
  media_type varchar NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  user_id bigint,
  CONSTRAINT message_log_primary_key PRIMARY KEY (id)

);

CREATE INDEX users_username_index ON aaa.users
  USING btree
  (
    username
  );

CREATE INDEX users_token_index ON aaa.users
  USING btree
  (
    token
  );

CREATE INDEX role_name_index ON aaa.roles
  USING btree
  (
    name
  );

ALTER TABLE aaa.user_roles ADD CONSTRAINT users_fk FOREIGN KEY (user_id)
REFERENCES aaa.users (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE aaa.user_roles ADD CONSTRAINT roles_fk FOREIGN KEY (role_id)
REFERENCES aaa.roles (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE aaa.message_logs ADD CONSTRAINT users_fk FOREIGN KEY (user_id)
REFERENCES aaa.users (id) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;

INSERT INTO aaa.roles (name, description, resources, created_at, updated_at) VALUES ('user', 'default role in system', '{}', 'now', 'now');
INSERT INTO aaa.roles (name, description, resources, created_at, updated_at) VALUES ('root', 'root group in system', '{god}', 'now', 'now');
-- Password is bita123
INSERT INTO aaa.users ("username", "password", "contact",  "token") VALUES ('admin', '$2a$10$6WeBOWQn2CwYzosiPK0ii.6XiW1rt0hZD3iXDsaySGo.RLoJUFwdq','admin@azmoona.com',  '251a1019f1370747b0093d5e9712124da88028f9');
-- Add all roles to admin user
INSERT INTO aaa.user_roles (user_id, role_id) SELECT (SELECT id FROM aaa.users WHERE username='admin') user_id , roles.id role_id FROM aaa.roles;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP SCHEMA aaa CASCADE;
