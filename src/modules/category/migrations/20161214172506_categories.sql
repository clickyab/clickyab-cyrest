
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE categories
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    scope VARCHAR(10) NOT NULL,
    title VARCHAR(30) NOT NULL,
    description VARCHAR(160),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE categories;
