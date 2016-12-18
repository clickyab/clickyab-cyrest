
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE users
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    email VARCHAR(50) NOT NULL UNIQUE ,
    password VARCHAR(60)NOT NULL ,
    old_password VARCHAR(30),
    access_token VARCHAR(60) NOT NULL ,
    user_type ENUM('personal', 'corporation'),
    parent_id INT(11),
    avatar VARCHAR(160),
    status ENUM('registered', 'verified', 'blocked') NOT NULL ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
    CONSTRAINT user_id_parent_id_fk FOREIGN KEY (parent_id) REFERENCES users (id)
);

CREATE INDEX users_users_id_fk ON users (parent_id);



CREATE TABLE roles
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(60) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL
);
CREATE UNIQUE INDEX role_id_uindex ON roles (id);


CREATE TABLE role_permission
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    role_id INT(11) NOT NULL,
    permission VARCHAR(255),
    scope ENUM('global', 'parent', 'self'),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
    CONSTRAINT role_permission_roles_id_fk FOREIGN KEY (role_id) REFERENCES roles (id)
);
CREATE UNIQUE INDEX role_permission_id_uindex ON role_permission (id);
CREATE INDEX role_permission_role_id_fk ON role_permission (role_id);

CREATE TABLE user_attributes
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id INT(11),
    `key` VARCHAR(100),
    value VARCHAR(100),
    CONSTRAINT user_attributes_users_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE UNIQUE INDEX user_attributes_id_uindex ON user_attributes (id);
CREATE INDEX user_attributes_users_id_fk ON user_attributes (user_id);

CREATE TABLE user_financial
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id INT(11) NOT NULL,
    bank_name VARCHAR(255),
    account_holder VARCHAR(255),
    card_number VARCHAR(2555),
    account_number VARCHAR(255),
    sheba_number VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
    CONSTRAINT user_financial_users_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE UNIQUE INDEX user_financial_id_uindex ON user_financial (id);
CREATE INDEX user_financial_user_id_fk ON user_financial (user_id);

CREATE TABLE user_profile_corporation
(
    user_id INT(11) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    economic_code VARCHAR(20),
    register_code VARCHAR(20) NOT NULL,
    phone VARCHAR(13),
    address VARCHAR(255),
    country_id INT(11),
    province_id INT(11),
    city_id INT(11),
    created_at DATETIME,
    updated_at DATETIME,
    CONSTRAINT user_profile_corporation_users_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX corporation_profile_city_id_fk ON user_profile_corporation (city_id);
CREATE INDEX corporation_profile_country_id_fk ON user_profile_corporation (country_id);
CREATE INDEX corporation_profile_province_id_fk ON user_profile_corporation (province_id);
CREATE UNIQUE INDEX corporation_profile_user_id_uindex ON user_profile_corporation (user_id);

CREATE TABLE user_profile_personal
(
    user_id INT(11) UNIQUE NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    birthday DATETIME,
    gender ENUM('male', 'female'),
    cellphone VARCHAR(255),
    phone VARCHAR(255),
    address TEXT,
    zip_code VARCHAR(255),
    national_code VARCHAR(10),
    country_id INT(11),
    province_id INT(11),
    city_id INT(11),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    CONSTRAINT user_profile_personal_users_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX personal_profile_city_id_fk ON user_profile_personal (city_id);
CREATE INDEX personal_profile_country_id_fk ON user_profile_personal (country_id);
CREATE INDEX personal_profile_province_id_fk ON user_profile_personal (province_id);
CREATE UNIQUE INDEX personal_profile_user_id_uindex ON user_profile_personal (user_id);


CREATE TABLE user_role
(
    user_id INT(11) DEFAULT '0' NOT NULL,
    role_id INT(11) DEFAULT '0' NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT user_role_users_id_fk FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT user_role_roles_id_fk FOREIGN KEY (role_id) REFERENCES roles (id)
);
CREATE INDEX user_role_role_id_fk ON user_role (role_id);
CREATE INDEX user_role_users_id_fk ON user_role (user_id);

INSERT INTO users (id,email,password,access_token,user_type,status,created_at,updated_at) VALUES (NULL,"root@rubik.com","$2a$10$6WeBOWQn2CwYzosiPK0ii.6XiW1rt0hZD3iXDsaySGo.RLoJUFwdq","92d80885abad94e24d3ffaea7501331fc7701135","personal","registered",NOW(),NOW());
INSERT INTO roles (id,name,description,created_at,updated_at) VALUES (NULL,"root","all access granted",NOW(),NOW());
INSERT INTO roles (id,name,description,created_at,updated_at) VALUES (NULL,"user","only some perm",NOW(),NOW());
INSERT INTO role_permission (id,role_id,permission,scope,created_at,updated_at) VALUES (NULL,(SELECT id FROM roles WHERE name="root"),"god","global",NOW(),NOW());
INSERT INTO user_role (user_id,role_id,created_at) VALUES ((SELECT id FROM users WHERE email="root@rubik.com"),(SELECT id FROM roles WHERE name="root"),NOW());
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE user_role;
DROP TABLE user_profile_personal;
DROP TABLE user_profile_corporation;
DROP TABLE user_financial;
DROP TABLE user_attributes;
DROP TABLE role_permission;
DROP TABLE roles;
DROP TABLE users;
