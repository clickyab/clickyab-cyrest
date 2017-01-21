
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT into plans (name,active,description,price,type,view) VALUES ("plan 1",1,"plan1plan1plan1plan1",100000,1,10000000);
INSERT into plans (name,active,description,price,type,view) VALUES ("plan 1",2,"plan1plan1plan1plan1",100000,1,10000000);
INSERT into plans (name,active,description,price,type,view) VALUES ("plan 1",1,"plan1plan1plan1plan1",100000,2,10000000);
INSERT into plans (name,active,description,price,type,view) VALUES ("plan 1",2,"plan1plan1plan1plan1",100000,2,10000000);
INSERT INTO ads (user_id,description,admin_status,archive_status,bot_chat_id,bot_message_id,cli_message_id,name,pay_status,position,plan_id) VALUES (1,"ad ad ad ad ad ",1,1,70018667,11,"05000000ce7365410900000000000000445450b6d9282f03","ads 1",1,10,1);
INSERT INTO ads (user_id,description,admin_status,archive_status,bot_chat_id,bot_message_id,cli_message_id,name,pay_status,position,plan_id) VALUES (1,"ad ad ad ad ad ",2,1,70018667,11,"05000000ce7365410900000000000000445450b6d9282f03","ads 1",1,10,2);
INSERT INTO ads (user_id,description,admin_status,archive_status,bot_chat_id,bot_message_id,cli_message_id,name,pay_status,position,plan_id) VALUES (1,"ad ad ad ad ad ",3,1,70018667,11,"05000000ce7365410900000000000000445450b6d9282f03","ads 1",1,10,3);
INSERT INTO ads (user_id,description,admin_status,archive_status,bot_chat_id,bot_message_id,cli_message_id,name,pay_status,position,plan_id) VALUES (1,"ad ad ad ad ad ",1,1,70018667,11,"05000000ce7365410900000000000000445450b6d9282f03","ads 1",1,10,1);
INSERT INTO ads (user_id,description,admin_status,archive_status,bot_chat_id,bot_message_id,cli_message_id,name,pay_status,position,plan_id) VALUES (1,"ad ad ad ad ad ",1,2,70018667,11,"05000000ce7365410900000000000000445450b6d9282f03","ads 1",1,10,2);
INSERT INTO ads (user_id,description,admin_status,archive_status,bot_chat_id,bot_message_id,cli_message_id,name,pay_status,position,plan_id) VALUES (1,"ad ad ad ad ad ",1,1,70018667,11,"05000000ce7365410900000000000000445450b6d9282f03","ads 1",1,10,3);
INSERT INTO ads (user_id,description,admin_status,archive_status,bot_chat_id,bot_message_id,cli_message_id,name,pay_status,position,plan_id) VALUES (1,"ad ad ad ad ad ",1,1,70018667,11,"05000000ce7365410900000000000000445450b6d9282f03","ads 1",2,10,1);
INSERT INTO users (id,email,password,access_token,user_type,status,created_at,updated_at) VALUES (NULL,"publisher@rubik.com","$2a$10$6WeBOWQn2CwYzosiPK0ii.6XiW1rt0hZD3iXDsaySGo.RLoJUFwdq","92d80885abad94edgfgd01331fc7701135","personal","registered",NOW(),NOW());
INSERT INTO users (id,email,password,access_token,user_type,status,created_at,updated_at) VALUES (NULL,"advertiser@rubik.com","$2a$10$6WeBOWQn2CwYzosiPK0ii.6XiW1rt0hZD3iXDsaySGo.RLoJUFwdq","92d80885abad9dfgdf31fc7701135","personal","registered",NOW(),NOW());
INSERT INTO telegram_users (user_id,bot_chat_id,remove,resolve,username) VALUES (1,49670863,2,1,"mahm0ud22");
INSERT INTO telegram_users (user_id,bot_chat_id,remove,resolve,username) VALUES (1,70018667,2,1,"AhmadDara");
INSERT INTO telegram_users (user_id,bot_chat_id,remove,resolve,username) VALUES (1,63205818,2,1,"mazafard");
INSERT INTO channels (user_id, name, link, admin_status, archive_status, active) VALUES (1,"chanel 1","telegram.me/1",3,2,1);
INSERT INTO channels (user_id, name, link, admin_status, archive_status, active) VALUES (1,"chanel 2","telegram.me/2",3,2,1);
INSERT INTO channels (user_id, name, link, admin_status, archive_status, active) VALUES (1,"chanel 3","telegram.me/3",3,2,1);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
 delete from users;
 delete from plans;
 delete from channels;
 delete from ads;
 delete from telegram_users;

