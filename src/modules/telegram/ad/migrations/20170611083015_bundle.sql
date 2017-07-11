
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE bundles
(
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    place INT,
    view INT NOT NULL DEFAULT 0,
    price INT NOT NULL,
    code VARCHAR(10) NOT NULL,
    percent_finish INT NOT NULL,
    bundle_type ENUM("banner", "banner+rep", "rep+banner", "rep+banner+rep") DEFAULT "banner" NOT NULL,
    rules TEXT,
    admin_status ENUM("on", "off") DEFAULT "off" NOT NULL,
    active_status ENUM("on", "off") DEFAULT "on" NOT NULL,
    ads VARCHAR(50) NOT NULL,
    target_ad INT NOT NULL,
    start TIMESTAMP,
    end TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT bundle_ads_id_fk FOREIGN KEY (target_ad) REFERENCES ads (id),
    CONSTRAINT bundle_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE bundles;