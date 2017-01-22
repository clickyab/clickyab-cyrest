
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE ads
(
  id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  user_id INT(11) NOT NULL,
  plan_id INT(11),
  name VARCHAR(60) NOT NULL,
  description TEXT,
  src TEXT,
  position INT(11),
  admin_status ENUM('pending', 'rejected','accepted'),
  archive_status ENUM('yes', 'no'),
  pay_status ENUM('yes', 'no'),
  active_status ENUM('yes','no'),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
  CONSTRAINT ads_id_users_id_fk FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT plans_id_plan_id_fk FOREIGN KEY (plan_id) REFERENCES plans(id)

);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE ads;

