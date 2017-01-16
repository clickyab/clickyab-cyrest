
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE files
(
  id INT(11) AUTO_INCREMENT PRIMARY KEY,
  user_id INT(11) NOT NULL,
  real_name VARCHAR(80) NOT NULL,
  db_name VARCHAR(80) NOT NULL,
  src TEXT,
  type ENUM('image','video','document'),
  size INT(11) UNSIGNED NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
  CONSTRAINT files_id_users_id_fk FOREIGN KEY (user_id) REFERENCES users(id)
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE files;

