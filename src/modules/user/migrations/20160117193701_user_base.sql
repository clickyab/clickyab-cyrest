
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE users
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    email VARCHAR(50) NOT NULL UNIQUE ,
    pub_name VARCHAR(50),
    password VARCHAR(60)NOT NULL ,
    old_password VARCHAR(30),
    access_token VARCHAR(60) NOT NULL ,
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


CREATE TABLE user_profile
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
    CONSTRAINT user_profile_users_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX profile_city_id_fk ON user_profile (city_id);
CREATE INDEX profile_country_id_fk ON user_profile (country_id);
CREATE INDEX profile_province_id_fk ON user_profile (province_id);
CREATE UNIQUE INDEX profile_user_id_uindex ON user_profile (user_id);


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

INSERT INTO users (id,email,password,access_token,status,created_at,updated_at) VALUES (NULL,"root@rubik.com","$2a$10$6WeBOWQn2CwYzosiPK0ii.6XiW1rt0hZD3iXDsaySGo.RLoJUFwdq","92d80885abad94e24d3ffaea7501331fc7701135","registered",NOW(),NOW());
INSERT INTO roles (id,name,description,created_at,updated_at) VALUES (1,"root","all access granted",NOW(),NOW());
INSERT INTO roles (id,name,description,created_at,updated_at) VALUES (2,"user","only some perm",NOW(),NOW());
INSERT INTO roles (id,name,description,created_at,updated_at) VALUES (3,"publisher","only publisher perm",NOW(),NOW());
INSERT INTO role_permission (id,role_id,permission,scope,created_at,updated_at) VALUES (NULL,(SELECT id FROM roles WHERE name="root"),"god","global",NOW(),NOW());
INSERT INTO user_role (user_id,role_id,created_at) VALUES ((SELECT id FROM users WHERE email="root@rubik.com"),(SELECT id FROM roles WHERE name="root"),NOW());
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE user_role;
DROP TABLE user_profile;
DROP TABLE user_financial;
DROP TABLE user_attributes;
DROP TABLE role_permission;
DROP TABLE roles;
DROP TABLE users;
