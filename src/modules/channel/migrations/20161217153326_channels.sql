
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE channels
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id INT(11)  NOT NULL,
    name VARCHAR(60) NOT NULL,
    link VARCHAR(100) ,
    admin VARCHAR(30)  ,
    status ENUM('pending', 'rejected','accepted','archive'),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
    CONSTRAINT channels_id_users_id_fk FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE channels;
