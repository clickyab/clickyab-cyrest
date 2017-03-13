
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE cyrest.billings MODIFY type ENUM('billing', 'withdrawal', 'income', 'campaign') DEFAULT 'billing';
ALTER TABLE cyrest.billings ADD channel_id INT NULL;
ALTER TABLE cyrest.billings ADD CONSTRAINT billings_channels_id_fk FOREIGN KEY (channel_id) REFERENCES channels (id);
ALTER TABLE cyrest.billings ADD ad_id INT NULL;
ALTER TABLE cyrest.billings ADD CONSTRAINT billings_ads_id_fk FOREIGN KEY (ad_id) REFERENCES ads (id);
ALTER TABLE cyrest.plans ADD share INT DEFAULT 40 NULL;

CREATE TABLE cyrest.billing_detail
(
    id INT PRIMARY KEY AUTO_INCREMENT,
    billing_id INT NOT NULL,
    user_id INT NOT NULL COMMENT 'who is doing work',
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT billing_detail_billings_id_fk FOREIGN KEY (billing_id) REFERENCES billings (id),
    CONSTRAINT billing_detail_users_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE cyrest.billing_detail;

