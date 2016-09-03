
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE users
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    email VARCHAR(255),
    password VARCHAR(255),
    old_password VARCHAR(255),
    access_token VARCHAR(255),
    source ENUM('crm', 'clickyab'),
    user_type ENUM('personal', 'corporation'),
    parent_id INT(11),
    avatar VARCHAR(255),
    status ENUM('registered', 'verified', 'blocked'),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL
);
CREATE TABLE user_role
(
    user_id INT(11) DEFAULT '0' NOT NULL,
    role_id INT(11) DEFAULT '0' NOT NULL,
    created_at TIMESTAMP,
    CONSTRAINT `PRIMARY` PRIMARY KEY (user_id, role_id)
);
CREATE TABLE user_profile_personal
(
    user_id INT(11) PRIMARY KEY NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    birthday DATETIME,
    gender ENUM('Male', 'Female'),
    cellphone VARCHAR(255),
    phone VARCHAR(255),
    address TEXT,
    zip_code VARCHAR(255),
    national_code VARCHAR(10),
    country_id INT(11),
    province_id INT(11),
    city_id INT(11),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
CREATE TABLE user_profile_corporation
(
    user_id INT(11) PRIMARY KEY NOT NULL,
    title VARCHAR(255),
    economic_code INT(11),
    register_code INT(11),
    phone INT(11),
    address VARCHAR(255),
    country_id INT(11),
    province_id INT(11),
    city_id INT(11),
    created_at DATETIME,
    update_at DATETIME
);
CREATE TABLE user_financial
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id INT(11) NOT NULL,
    bank_name VARCHAR(255),
    account_holder VARCHAR(255),
    card_number VARCHAR(2555),
    account_number VARCHAR(255),
    sheba_number VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE TABLE user_crm
(
    user_id INT(11) PRIMARY KEY NOT NULL,
    originating_lead VARCHAR(255),
    customer_code INT(12),
    gid VARCHAR(64),
    lead_status INT(1),
    read_by_crm INT(2)
);
CREATE TABLE user_attributes
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id INT(11),
    `key` VARCHAR(100),
    value VARCHAR(100)
);
CREATE TABLE roles
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL
);
CREATE TABLE role_permission
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    role_id INT(11) NOT NULL,
    permission VARCHAR(255),
    scope ENUM('global', 'parent', 'own'),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL
);
CREATE TABLE provinces
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    country_id INT(11) NOT NULL,
    name_en VARCHAR(255),
    name_fa VARCHAR(255),
    region_code INT(2),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE TABLE countries
(
    id INT(11) PRIMARY KEY NOT NULL,
    iso VARCHAR(2) NOT NULL,
    name VARCHAR(100) NOT NULL,
    iso3 VARCHAR(3),
    country_code SMALLINT(6),
    phone_code INT(5),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE TABLE cities
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    province_id INT(11) NOT NULL,
    name_en VARCHAR(255),
    name_fa VARCHAR(255),
    city_code INT(11),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE TABLE message_log
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id INT(11),
    subject VARCHAR(255),
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL,
    status ENUM('read', 'pending', 'sent'),
    CONSTRAINT message_log_users_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);
ALTER TABLE users ADD FOREIGN KEY (parent_id) REFERENCES users (id);
CREATE INDEX users_users_id_fk ON users (parent_id);
CREATE UNIQUE INDEX user_id_uindex ON users (id);
ALTER TABLE user_role ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE user_role ADD FOREIGN KEY (role_id) REFERENCES roles (id);
CREATE INDEX user_role_role_id_fk ON user_role (role_id);
ALTER TABLE user_profile_personal ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE user_profile_personal ADD FOREIGN KEY (country_id) REFERENCES countries (id);
ALTER TABLE user_profile_personal ADD FOREIGN KEY (province_id) REFERENCES provinces (id);
ALTER TABLE user_profile_personal ADD FOREIGN KEY (city_id) REFERENCES cities (id);
CREATE INDEX personal_profile_city_id_fk ON user_profile_personal (city_id);
CREATE INDEX personal_profile_country_id_fk ON user_profile_personal (country_id);
CREATE INDEX personal_profile_province_id_fk ON user_profile_personal (province_id);
CREATE UNIQUE INDEX personal_profile_user_id_uindex ON user_profile_personal (user_id);
ALTER TABLE user_profile_corporation ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE user_profile_corporation ADD FOREIGN KEY (country_id) REFERENCES countries (id);
ALTER TABLE user_profile_corporation ADD FOREIGN KEY (province_id) REFERENCES provinces (id);
ALTER TABLE user_profile_corporation ADD FOREIGN KEY (city_id) REFERENCES cities (id);
CREATE INDEX corporation_profile_city_id_fk ON user_profile_corporation (city_id);
CREATE INDEX corporation_profile_country_id_fk ON user_profile_corporation (country_id);
CREATE INDEX corporation_profile_province_id_fk ON user_profile_corporation (province_id);
CREATE UNIQUE INDEX corporation_profile_user_id_uindex ON user_profile_corporation (user_id);
ALTER TABLE user_financial ADD FOREIGN KEY (user_id) REFERENCES users (id);
CREATE UNIQUE INDEX user_financial_id_uindex ON user_financial (id);
CREATE INDEX user_financial_user_id_fk ON user_financial (user_id);
ALTER TABLE user_crm ADD FOREIGN KEY (user_id) REFERENCES users (id);
CREATE UNIQUE INDEX user_crm_user_id_uindex ON user_crm (user_id);
ALTER TABLE user_attributes ADD FOREIGN KEY (user_id) REFERENCES users (id);
CREATE UNIQUE INDEX user_attributes_id_uindex ON user_attributes (id);
CREATE INDEX user_attributes_users_id_fk ON user_attributes (user_id);
CREATE UNIQUE INDEX role_id_uindex ON roles (id);
ALTER TABLE role_permission ADD FOREIGN KEY (role_id) REFERENCES roles (id);
CREATE UNIQUE INDEX role_permission_id_uindex ON role_permission (id);
CREATE INDEX role_permission_role_id_fk ON role_permission (role_id);
ALTER TABLE provinces ADD FOREIGN KEY (country_id) REFERENCES countries (id);
CREATE INDEX province_country_id_fk ON provinces (country_id);
CREATE UNIQUE INDEX province_id_uindex ON provinces (id);
CREATE UNIQUE INDEX Country_iso_uindex ON countries (iso);
CREATE UNIQUE INDEX Country_name_uindex ON countries (name);
ALTER TABLE cities ADD FOREIGN KEY (province_id) REFERENCES provinces (id);
CREATE UNIQUE INDEX city_id_uindex ON cities (id);
CREATE INDEX city_province_id_fk ON cities (province_id);
CREATE INDEX message_log_users_id_fk ON message_log (user_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE user_role;
DROP TABLE message_log;
DROP TABLE user_profile_personal;
DROP TABLE user_profile_corporation;
DROP TABLE user_financial;
DROP TABLE user_crm;
DROP TABLE user_attributes;
DROP TABLE role_permission;
DROP TABLE roles;
DROP TABLE cities;
DROP TABLE provinces;
DROP TABLE countries;
DROP TABLE users;
