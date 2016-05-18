
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE acc.accounts DROP COLUMN "initial";

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE acc.accounts ADD COLUMN "initial" bigint NOT NULL DEFAULT 0;