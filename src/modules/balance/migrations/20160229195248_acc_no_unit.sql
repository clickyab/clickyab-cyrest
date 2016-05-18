
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

DROP SCHEMA IF EXISTS acc CASCADE;

CREATE SCHEMA acc;

CREATE TABLE acc.accounts(
	id bigserial NOT NULL,
	owner_id bigint NOT NULL,
	title varchar NOT NULL,
	description text,
	created_at timestamp with time zone NOT NULL DEFAULT NOW(),
	updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
	attributes jsonb NOT NULL DEFAULT '{}',
	disabled bool NOT NULL DEFAULT false,
	CONSTRAINT accounts_primary_keys PRIMARY KEY (id)

);

CREATE TABLE acc.transactions(
	id bigserial NOT NULL,
	account_id bigint NOT NULL,
	user_id bigint NOT NULL,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone NOT NULL,
	amount bigint NOT NULL,
	description text,
	correlation_id varchar DEFAULT NULL,
	CONSTRAINT transactions_primary_key PRIMARY KEY (id),
	CONSTRAINT amount_must_not_be_zero CHECK (amount <> 0)

);

CREATE INDEX transaction_account_index ON acc.transactions
	USING btree
	(
	  account_id
	);

CREATE INDEX transaction_user_id_foreign ON acc.transactions
	USING btree
	(
	  user_id
	);

CREATE TABLE acc.transaction_tags(
	transaction_id bigint,
	tags varchar[],
	CONSTRAINT transaction_tags_pk PRIMARY KEY (transaction_id)

);

CREATE INDEX transaction_tags_index_gin ON acc.transaction_tags
	USING gin
	(
	  tags
	);

CREATE TABLE acc.user_tags(
	id bigserial NOT NULL,
	tag text NOT NULL,
	count integer DEFAULT 0,
	user_id bigint NOT NULL
);

ALTER TABLE acc.transactions ADD CONSTRAINT accounts_fk FOREIGN KEY (account_id)
REFERENCES acc.accounts (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE acc.transactions ADD CONSTRAINT users_fk FOREIGN KEY (user_id)
REFERENCES aaa.users (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE acc.transaction_tags ADD CONSTRAINT transactions_fk FOREIGN KEY (transaction_id)
REFERENCES acc.transactions (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE acc.user_tags ADD CONSTRAINT users_fk FOREIGN KEY (user_id)
REFERENCES aaa.users (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE acc.accounts ADD CONSTRAINT users_fk FOREIGN KEY (owner_id)
REFERENCES aaa.users (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP SCHEMA acc;