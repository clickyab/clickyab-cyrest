
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE ads

  ADD COLUMN `cli_message_id` VARCHAR (70) AFTER `src`,
  ADD COLUMN `bot_chat_id` INT,
  ADD COLUMN `bot_message_id` INT,
  MODIFY COLUMN src VARCHAR(255);
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE ads
  DROP COLUMN `cli_message_id`,
  DROP COLUMN `bot_chat_id`,
  DROP COLUMN `bot_message_id`,
  MODIFY COLUMN src TEXT;

