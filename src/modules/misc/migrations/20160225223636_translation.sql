
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE SCHEMA t9n;

CREATE TABLE t9n.translations(
	id bigserial NOT NULL,
	string text NOT NULL,
	single text DEFAULT null,
	plural text DEFAULT null,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone NOT NULL,
	CONSTRAINT translations_primary_key PRIMARY KEY (id),
	CONSTRAINT unique_string_translation UNIQUE (string)
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP SCHEMA t9n CASCADE;
