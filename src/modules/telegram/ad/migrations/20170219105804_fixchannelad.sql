
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE channel_ad DROP COLUMN start;
ALTER TABLE channel_ad DROP COLUMN end;
ALTER TABLE channel_ad DROP ADD COLUMN start TIMESTAMP NULL NULL;
ALTER TABLE channel_ad DROP ADD COLUMN end TIMESTAMP NULL NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


