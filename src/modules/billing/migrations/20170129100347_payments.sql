
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE payments
(
  id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  user_id INT(11)  NOT NULL,
  amount INT(11) NOT NULL,
  status ENUM('pending','rejected','paid') NOT NULL,
  authority VARCHAR(100),
  ref_id  INT(11),
  status_code  INT(11),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
  CONSTRAINT payments._id_users_id_fk FOREIGN KEY (user_id) REFERENCES users(id)
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE payments;