
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE plans
  ADD COLUMN price INT(11) NOT NULL AFTER description;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE plans
  DROP COLUMN price;

