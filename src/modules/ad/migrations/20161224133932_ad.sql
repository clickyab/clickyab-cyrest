
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE ads
(
  id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  user_id INT(11) NOT NULL,
  name VARCHAR(60) NOT NULL,
  type ENUM('img','document','video'),
  media TEXT,
  description TEXT,
  link TEXT,
  status ENUM('pending', 'rejected','accepted','archive'),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
  CONSTRAINT ads_id_users_id_fk FOREIGN KEY (user_id) REFERENCES users(id)

);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE ads;

