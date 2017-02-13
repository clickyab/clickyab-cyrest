
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE UNIQUE INDEX channel_details_channel_id_uindex ON cyrest.channel_details (channel_id);
DROP INDEX name ON cyrest.channel_details;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


