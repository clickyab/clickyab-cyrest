
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE cyrest.billings ADD type ENUM("billing", "withdrawal") DEFAULT "billing" NULL;
ALTER TABLE cyrest.billings ADD status ENUM("accepted", "rejected", "pending") DEFAULT "pending" NULL;
ALTER TABLE cyrest.billings ADD deposit ENUM("yes", "no") DEFAULT "no" NULL;
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


