
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE ad_category
(
    ad_id INT NOT NULL,
    category_id INT NOT NULL,
    CONSTRAINT ad_category_ads_id_fk FOREIGN KEY (ad_id) REFERENCES ads (id),
    CONSTRAINT ad_category_categories_id_fk FOREIGN KEY (category_id) REFERENCES categories (id)
);
CREATE UNIQUE INDEX ad_category_ad_id_category_id_uindex ON cyrest.ad_category (ad_id, category_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE ad_category;

