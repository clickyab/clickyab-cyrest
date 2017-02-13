
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE UNIQUE INDEX translations_string_id_lang_uindex ON cyrest.translations (string_id, lang);
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


