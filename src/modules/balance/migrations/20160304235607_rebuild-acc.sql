
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE acc.accounts ADD COLUMN current bigint NOT NULL DEFAULT 0;
ALTER TABLE acc.user_tags ADD CONSTRAINT user_tag_uniq UNIQUE (user_id,tag);

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
            INSERT INTO acc.user_tags (user_id, tag, count)
                SELECT (SELECT owner_id FROM acc.accounts a LEFT JOIN acc.transactions t ON a.id=t.account_id
                WHERE t.id = NEW.transaction_id),tag, 1 FROM UNNEST(NEW.tags) tag
                ON CONFLICT ON CONSTRAINT user_tag_uniq DO UPDATE SET count = acc.user_tags.count +1;
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

ALTER TABLE acc.accounts DROP COLUMN current;
ALTER TABLE acc.user_tags DROP CONSTRAINT user_tag_uniq;
DROP TRIGGER insert_transaction ON acc.transactions;
DROP TRIGGER insert_transaction_tags ON acc.transaction_tags;
DROP FUNCTION acc.insert_transaction_tags();
DROP FUNCTION acc.insert_transaction();