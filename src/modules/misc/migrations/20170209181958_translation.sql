
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

DROP TABLE IF EXISTS translations CASCADE;

CREATE TABLE strings(
        id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
        text TEXT NOT NULL
);

CREATE TABLE translations(
       id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
       string_id INT(11) NOT NULL,
       lang  VARCHAR(10) NOT NULL,
       translated TEXT NOT NULL,
       CONSTRAINT translations_string_id_fk FOREIGN KEY (string_id) REFERENCES strings (id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


DROP TABLE IF EXISTS translations CASCADE;
DROP TABLE IF EXISTS strings CASCADE;

CREATE TABLE translations
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    string TEXT NOT NULL,
    single TEXT,
    plural TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL
);