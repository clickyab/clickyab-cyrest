
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied


CREATE TABLE plans
(
  id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  name VARCHAR(60) NOT NULL,
  description TEXT,
  price INT(11),
  view INT(11)

);

CREATE TABLE ads
(
  id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  user_id INT(11) NOT NULL,
  plan_id INT(11),
  name VARCHAR(60) NOT NULL,
  description TEXT,
  src TEXT,
  admin_status ENUM('pending', 'rejected','accepted'),
  archive_status ENUM('yes', 'no'),
  pay_status ENUM('yes', 'no'),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
  CONSTRAINT ads_id_users_id_fk FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT plans_id_plan_id_fk FOREIGN KEY (plan_id) REFERENCES plans(id)

);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE plans;
DROP TABLE ads;

