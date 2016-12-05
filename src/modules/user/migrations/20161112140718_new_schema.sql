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
CREATE INDEX users_users_id_fk ON users (parent_id);
CREATE UNIQUE INDEX user_id_uindex ON users (id);

CREATE TABLE roles
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT '0000-00-00 00:00:00' NOT NULL
);
CREATE UNIQUE INDEX role_id_uindex ON roles (id);


CREATE TABLE role_permission
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    role_id INT(11) NOT NULL,
    permission VARCHAR(255),
    scope ENUM('global', 'parent', 'own'),
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
CREATE TABLE user_crm
(
    user_id INT(11) NOT NULL,
    originating_lead VARCHAR(255),
    customer_code INT(12),
    gid VARCHAR(64),
    lead_status INT(1),
    read_by_crm INT(2),
    CONSTRAINT user_crm_users_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE UNIQUE INDEX user_crm_user_id_uindex ON user_crm (user_id);

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
    user_id INT(11) NOT NULL,
    title VARCHAR(255),
    economic_code INT(11),
    register_code INT(11),
    phone INT(11),
    address VARCHAR(255),
    country_id INT(11),
    province_id INT(11),
    city_id INT(11),
    created_at DATETIME,
    update_at DATETIME,
    CONSTRAINT user_profile_corporation_users_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX corporation_profile_city_id_fk ON user_profile_corporation (city_id);
CREATE INDEX corporation_profile_country_id_fk ON user_profile_corporation (country_id);
CREATE INDEX corporation_profile_province_id_fk ON user_profile_corporation (province_id);
CREATE UNIQUE INDEX corporation_profile_user_id_uindex ON user_profile_corporation (user_id);

CREATE TABLE user_profile_personal
(
    user_id INT(11) NOT NULL,
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
