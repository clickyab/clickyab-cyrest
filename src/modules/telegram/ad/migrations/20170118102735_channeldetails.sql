
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE channel_details
(
  id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(60) NOT NULL UNIQUE,
  channel_id INT(11) NOT NULL,
  title VARCHAR(60)NOT NULL,
  info TEXT,
  cli_telegram_id VARCHAR(60),
  user_count INT(11),
  admin_count INT(11),
  post_count INT(11),
  total_view INT(11),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
  CONSTRAINT channel_details_channels_id_fk FOREIGN KEY (channel_id) REFERENCES channels(id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE channel_details;

