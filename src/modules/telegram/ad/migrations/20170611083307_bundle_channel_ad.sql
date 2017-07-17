
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE bundle_channel_ad
(
    channel_id INT NOT NULL,
    bundle_id INT NOT NULL,
    ad_id INT NOT NULL,
    view INT DEFAULT 0,
    shot VARCHAR(255),
    warning TINYINT DEFAULT 0,
    bot_message_id INT(11),
    bot_chat_id INT(11) NOT NULL,
    active ENUM('yes','no') NOT NULL,
    start TIMESTAMP DEFAULT NULL,
    end TIMESTAMP DEFAULT NULL,
    shot_time TIMESTAMP DEFAULT NULL,
    billing_time TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (channel_id,ad_id,bundle_id),
    CONSTRAINT bundle_channel_ad_channel_id_fk FOREIGN KEY (channel_id) REFERENCES channels (id),
    CONSTRAINT bundle_channel_ad_bundle_id_fk FOREIGN KEY (bundle_id) REFERENCES bundles (id),
    CONSTRAINT bundle_channel_ad_ad_id_fk FOREIGN KEY (ad_id) REFERENCES ads (id)
);

CREATE INDEX channel_bundle_ad_channel_id_fk ON bundle_channel_ad (channel_id);
CREATE INDEX channel_bundle_ad_bundle_id_fk ON bundle_channel_ad (bundle_id);
CREATE INDEX channel_bundle_ad_ad_id_fk ON bundle_channel_ad (ad_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE bundle_channel_ad;

