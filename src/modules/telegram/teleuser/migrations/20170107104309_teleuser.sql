
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE telegram_users
(
  id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  user_id INT(11)  NOT NULL,
  username VARCHAR(60) NOT NULL,
  bot_chat_id INT UNIQUE,
  resolve ENUM('yes','no') NOT NULL,
  remove ENUM('yes','no') NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
  CONSTRAINT telegram_users._id_users_id_fk FOREIGN KEY (user_id) REFERENCES users(id)

);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE telegram_users;
