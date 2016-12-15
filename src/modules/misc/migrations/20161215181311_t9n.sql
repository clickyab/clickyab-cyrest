
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE translations
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    string TEXT NOT NULL,
    single TEXT,
    plural TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE translations;