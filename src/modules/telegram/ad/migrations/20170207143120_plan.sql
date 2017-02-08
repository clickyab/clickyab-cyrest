
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO plans (name, description, price, view, position, type, active, created_at, updated_at) VALUES ("promotion","promotion package",4000,1000000,5,"promotion","yes",NOW(),NOW());
INSERT INTO plans (name, description, price, view, position, type, active, created_at, updated_at) VALUES ("individual","individual package",3500,1000000,5,"individual","yes",NOW(),NOW());


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
truncate plans;

