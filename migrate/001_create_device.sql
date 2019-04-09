-- +migrate Up
create table device(
    device_eui bytea primary key,
    device_name varchar(200) not null,
    icon varchar(100) default '', 
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

-- +migrate Down
drop table device;