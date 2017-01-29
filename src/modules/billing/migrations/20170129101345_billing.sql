
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE billings
(
  id INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id INT(11)  NOT NULL,
  amount INT(11) NOT NULL,
  reason VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
  CONSTRAINT billing._id_users_id_fk FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE billings;

