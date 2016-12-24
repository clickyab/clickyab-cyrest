
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE campaigns
(
  id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  user_id INT(11)  NOT NULL,
  name VARCHAR(60) NOT NULL,
  active ENUM('yes','no'),
  start TIMESTAMP,
  end TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
  CONSTRAINT campaigns_id_users_id_fk FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE campaign_black
(
  campaign_id INT(11) NOT NULL,
  channel_id INT(11) NOT NULL,
  CONSTRAINT campaign_channel_channel_id_fk FOREIGN KEY (channel_id) REFERENCES channels (id),
  CONSTRAINT campaign_channel_campaign_id_id_fk FOREIGN KEY (campaign_id) REFERENCES campaigns (id)
);
CREATE INDEX campaign_channel_channel_id_fk ON campaign_black (channel_id);
CREATE INDEX campaign_channel_campaign_id_id_fk ON campaign_black (campaign_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE campaign_black;
DROP TABLE campaigns;


