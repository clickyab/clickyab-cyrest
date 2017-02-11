
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE ads

  ADD COLUMN `view` INT (11) DEFAULT 0 AFTER `bot_message_id`;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE ads
  DROP COLUMN `view`;
