
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE ads ADD COLUMN promo_src TEXT AFTER src;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE ads DROP COLUMN promo_src;


