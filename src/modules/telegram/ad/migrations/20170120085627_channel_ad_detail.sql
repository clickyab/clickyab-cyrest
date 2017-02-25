
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE channel_ad_detail
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    channel_id INT NOT NULL,
    ad_id INT NOT NULL,
    view INT,
    position INT,
    warning INT(1),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT channel_ad_detail_channel_id_fk FOREIGN KEY (channel_id) REFERENCES channels (id),
    CONSTRAINT channel_ad_detail_ad_id_fk FOREIGN KEY (ad_id) REFERENCES ads (id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE channel_ad_detail;
