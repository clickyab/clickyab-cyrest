
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE bundle_channel_ad_detail
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    channel_id INT NOT NULL,
    bundle_id INT NOT NULL,
    ad_id INT NOT NULL,
    view INT DEFAULT 0,
    position INT,
    warning INT(1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT bundle_channel_ad_detail_channel_id_fk FOREIGN KEY (channel_id) REFERENCES channels (id),
    CONSTRAINT bundle_channel_ad_detail_bundle_id_fk FOREIGN KEY (bundle_id) REFERENCES bundles (id),
    CONSTRAINT bundle_channel_ad_detail_ad_id_fk FOREIGN KEY (ad_id) REFERENCES ads (id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE bundle_channel_ad_detail;

