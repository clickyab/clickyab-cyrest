
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE plans
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(60) NOT NULL,
    description  TEXT,
    view INT,
    type ENUM('promotion','individual') NOT NULL,
    active ENUM('yes', 'no') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE plans;