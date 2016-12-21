
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

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE campaigns;


