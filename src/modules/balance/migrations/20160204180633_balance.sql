
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE SCHEMA acc;

CREATE TABLE acc.units(
  owner_id bigint,
  title varchar NOT NULL,
  description text,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  attributes jsonb NOT NULL DEFAULT '{}',
  CONSTRAINT units_pk PRIMARY KEY (owner_id)

);

CREATE TABLE acc.accounts(
  id bigserial NOT NULL,
  unit_id bigint,
  title varchar NOT NULL,
  description text,
  created_at timestamp with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
  initial bigint NOT NULL DEFAULT 0,
  attributes jsonb NOT NULL DEFAULT '{}',
  current bigint NOT NULL DEFAULT 0,
  disabled bool NOT NULL DEFAULT false,
  CONSTRAINT accounts_primary_keys PRIMARY KEY (id)

);

CREATE TABLE acc.transactions(
  id bigserial NOT NULL,
  created_at timestamp with time zone NOT NULL,
  updated_at timestamp with time zone NOT NULL,
  amount bigint NOT NULL,
  description text,
  account_id bigint NOT NULL,
  user_id bigint NOT NULL,
  CONSTRAINT transactions_primary_key PRIMARY KEY (id),
  CONSTRAINT amount_must_not_be_zero CHECK (amount <> 0)

);

CREATE TABLE acc.transaction_tags(
  transaction_id bigint,
  tags varchar[],
  CONSTRAINT transaction_tags_pk PRIMARY KEY (transaction_id)

);

CREATE TABLE acc.unit_tags(
	id bigserial NOT NULL,
	tag text NOT NULL,
	count integer DEFAULT 0,
	unit_id bigint NOT NULL
);
-- ddl-end --

-- object: units_fk | type: CONSTRAINT --
-- ALTER TABLE acc.unit_tags DROP CONSTRAINT IF EXISTS units_fk CASCADE;
ALTER TABLE acc.unit_tags ADD CONSTRAINT units_fk FOREIGN KEY (unit_id)
REFERENCES acc.units (owner_id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: unit_tag_uq | type: CONSTRAINT --
-- ALTER TABLE acc.unit_tags DROP CONSTRAINT IF EXISTS unit_tag_uq CASCADE;
ALTER TABLE acc.unit_tags ADD CONSTRAINT unit_tag_uq UNIQUE (tag,unit_id);


CREATE TABLE acc.unit_users(
  unit_id bigint,
  user_id bigint,
  permissions jsonb NOT NULL DEFAULT '{}',
  CONSTRAINT unit_users_pk PRIMARY KEY (unit_id,user_id)

);

ALTER TABLE acc.unit_users ADD CONSTRAINT units_fk FOREIGN KEY (unit_id)
REFERENCES acc.units (owner_id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE acc.unit_users ADD CONSTRAINT users_fk FOREIGN KEY (user_id)
REFERENCES aaa.users (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;

CREATE INDEX account_unit_id_foreign ON acc.accounts
USING btree
(
  unit_id
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

CREATE INDEX transaction_tags_index_gin ON acc.transaction_tags
USING gin
(
  tags
);

ALTER TABLE acc.units ADD CONSTRAINT users_fk FOREIGN KEY (owner_id)
REFERENCES aaa.users (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE acc.accounts ADD CONSTRAINT units_fk FOREIGN KEY (unit_id)
REFERENCES acc.units (owner_id) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE acc.transactions ADD CONSTRAINT accounts_fk FOREIGN KEY (account_id)
REFERENCES acc.accounts (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE acc.transactions ADD CONSTRAINT users_fk FOREIGN KEY (user_id)
REFERENCES aaa.users (id) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE acc.transaction_tags ADD CONSTRAINT transactions_fk FOREIGN KEY (transaction_id)
REFERENCES acc.transactions (id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;

INSERT INTO acc.units (owner_id, title) SELECT id, 'default home' FROM aaa.users;
INSERT INTO acc.accounts (unit_id, title) SELECT owner_id, 'personal pocket' FROM acc.units;

-- +migrate StatementBegin

CREATE FUNCTION acc.insert_transaction() RETURNS trigger AS $insert_transaction$
    BEGIN
        IF TG_OP = 'INSERT' THEN
           UPDATE acc.accounts SET "current"="current"+NEW.amount WHERE id=NEW.account_id;
        END IF;

        IF TG_OP = 'UPDATE' THEN
           UPDATE acc.accounts SET "current"="current"+NEW.amount-OLD.amount WHERE id=NEW.account_id;
        END IF;

        IF TG_OP = 'DELETE' THEN
           UPDATE acc.accounts SET "current"="current"-OLD.amount WHERE id=NEW.account_id;
        END IF;

        IF TG_OP = 'TRUNCATE' THEN
           UPDATE acc.accounts SET "current"=0;
        END IF;

        -- Ready to continue!
        RETURN NEW;
    END;
$insert_transaction$ LANGUAGE plpgsql;

CREATE TRIGGER insert_transaction BEFORE INSERT OR UPDATE OR DELETE ON acc.transactions
    FOR EACH ROW EXECUTE PROCEDURE acc.insert_transaction();


CREATE FUNCTION acc.insert_transaction_tags() RETURNS trigger AS $insert_transaction_tags$
    BEGIN
        IF TG_OP = 'INSERT' THEN
            INSERT INTO acc.unit_tags (unit_id, tag, count)
                SELECT (SELECT unit_id FROM acc.accounts a LEFT JOIN acc.transactions t ON a.id=t.account_id
                WHERE t.id = NEW.transaction_id),tag, 1 FROM UNNEST(NEW.tags) tag
                ON CONFLICT ON CONSTRAINT unit_tag_uq DO UPDATE SET count = acc.unit_tags.count +1;
        END IF;

        -- Ready to continue!
        RETURN NEW;
    END;
$insert_transaction_tags$ LANGUAGE plpgsql;

CREATE TRIGGER insert_transaction_tags BEFORE INSERT OR UPDATE OR DELETE ON acc.transaction_tags
    FOR EACH ROW EXECUTE PROCEDURE acc.insert_transaction_tags();

-- +migrate StatementEnd

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP SCHEMA acc CASCADE;
