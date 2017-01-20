
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE known_channels
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE ,
    title VARCHAR(60)NOT NULL ,
    info VARCHAR(160),
    user_count INT(11),
    telegram_id VARCHAR(60),
    raw_data TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL

);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE known_channels;