
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT into plans (name,active,description,price,type,view) VALUES ("پلن ۱","yes","توضیحات پلن ۱",100000,"individual",10);
INSERT into plans (name,active,description,price,type,view) VALUES ("پلن ۲","yes","توضیحات پلن ۲",120000,"promotion",100);
INSERT into plans (name,active,description,price,type,view) VALUES ("پلن ۳","yes","توضیحات پلن ۳",15000,"individual",1000);
INSERT into plans (name,active,description,price,type,view) VALUES ("پلن ۴","yes","توضیحات پلن ۴",250000,"promotion",1000);
INSERT INTO ads (user_id, plan_id, name, description, cli_message_id, position, admin_status, archive_status, pay_status, bot_chat_id, bot_message_id) VALUES (1,1,"ad ad ad ad ad ","testtesttesttesttesttesttest","05000000ce7365410900000000000000445450b6d9282f03",10,"accepted","no","yes",49670863,10);
INSERT INTO ads (user_id, plan_id, name, description, cli_message_id, position, admin_status, archive_status, pay_status, bot_chat_id, bot_message_id) VALUES (1,1,"ad ad ad ad ad ","testtesttesttesttesttesttest","05000000ce7365410900000000000000445450b6d9282f03",10,"accepted","no","yes",49670863,10);
INSERT INTO users (id,email,password,access_token,user_type,status,created_at,updated_at) VALUES (NULL,"publisher@rubik.com","$2a$10$6WeBOWQn2CwYzosiPK0ii.6XiW1rt0hZD3iXDsaySGo.RLoJUFwdq","92d80885abad94edgfgd01331fc7701135","personal","registered",NOW(),NOW());
INSERT INTO users (id,email,password,access_token,user_type,status,created_at,updated_at) VALUES (NULL,"advertiser@rubik.com","$2a$10$6WeBOWQn2CwYzosiPK0ii.6XiW1rt0hZD3iXDsaySGo.RLoJUFwdq","92d80885abad9dfgdf31fc7701135","personal","registered",NOW(),NOW());
INSERT INTO telegram_users (user_id,bot_chat_id,remove,resolve,username) VALUES (1,49670863,"no","yes","mahm0ud22");
INSERT INTO telegram_users (user_id,bot_chat_id,remove,resolve,username) VALUES (1,70018667,"no","yes","AhmadDara");
INSERT INTO telegram_users (user_id,bot_chat_id,remove,resolve,username) VALUES (1,63205818,"no","yes","mazafard");
INSERT INTO channels (user_id, name, link, admin_status, archive_status, active) VALUES (1,"test234","https://t.me/tst1234567","accepted","no","yes");
INSERT INTO channels (user_id, name, link, admin_status, archive_status, active) VALUES (1,"daratest","https://t.me/daratest","accepted","no","yes");
INSERT INTO channels (user_id, name, link, admin_status, archive_status, active) VALUES (1,"mamedagha","https://t.me/mamaedagha","accepted","no","yes");
INSERT channel_ad (channel_id, ad_id, view, cli_message_id, warning, active, start, possible_view) VALUES (1,1,10,"05000000ce7365410900000000000000445450b6d9282f03,1","yes",NOW(),200);
INSERT channel_ad (channel_id, ad_id, view, cli_message_id, warning, active, start, possible_view) VALUES (1,1,100,"05000000ce7365410900000000000000445450b6d9282f03,2","yes",NOW(),2000);
INSERT channel_ad (channel_id, ad_id, view, cli_message_id, warning, active, start, possible_view) VALUES (1,1,1000,"05000000ce7365410900000000000000445450b6d9282f03",4,"yes",NOW(),400);
INSERT channel_ad (channel_id, ad_id, view, cli_message_id, warning, active, start, possible_view) VALUES (1,1,10000,"05000000ce7365410900000000000000445450b6d9282f03",1,"yes",NOW(),200);
INSERT channel_ad (channel_id, ad_id, view, cli_message_id, warning, active, start, possible_view) VALUES (1,1,1000000,"05000000ce7365410900000000000000445450b6d9282f03",7,"yes",NOW(),200);
INSERT channel_ad (channel_id, ad_id, view, cli_message_id, warning, active, start, possible_view) VALUES (1,1,10000000,"05000000ce7365410900000000000000445450b6d9282f03",8,"yes",NOW(),200);
INSERT channel_ad (channel_id, ad_id, view, cli_message_id, warning, active, start, possible_view) VALUES (1,1,1000000,"05000000ce7365410900000000000000445450b6d9282f03",1,"yes",NOW(),200);
INSERT channel_ad (channel_id, ad_id, view, cli_message_id, warning, active, start, possible_view) VALUES (1,1,1000,"05000000ce7365410900000000000000445450b6d9282f03",1,"yes",NOW(),200);
INSERT channel_ad (channel_id, ad_id, view, cli_message_id, warning, active, start, possible_view) VALUES (1,1,1,"05000000ce7365410900000000000000445450b6d9282f03",1,"yes",NOW(),200);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
 delete from plans;
 delete from telegram_users;
 delete from users;
 delete from ads;
 delete from channels;
 delete from channel_ad;
 delete from channel_details;

