
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE channel_ad
(
    channel_id INT NOT NULL,
    ad_id INT NOT NULL,
    view INT,
    cli_message_id VARCHAR(70),
    warning INT,
    active ENUM('yes','no'),
    start TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
    end TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
    possible_view INT(11),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
    PRIMARY KEY (channel_id,ad_id),
    CONSTRAINT channel_ad_channel_id_fk FOREIGN KEY (channel_id) REFERENCES channels (id),
    CONSTRAINT channel_ad_ad_id_fk FOREIGN KEY (ad_id) REFERENCES ads (id)
);

CREATE INDEX channel_ad_channel_id_fk ON channel_ad (channel_id);
CREATE INDEX channel_ad_ad_id_fk ON channel_ad (ad_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE channel_ad;

