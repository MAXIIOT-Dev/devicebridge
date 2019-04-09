-- +migrate Up
create table users(
    user_name varchar(100) not null primary key,
    password_hash varchar(200) not null
);

insert into users(user_name,password_hash)
values('admin','$2a$10$cCLGdc9rmnwTkKdeR6LpHeniqp2ZvI9q6fWC7LDaKUz7dFcnrKBdi');
-- +migrate Down
drop table users;