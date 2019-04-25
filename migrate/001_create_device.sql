-- +migrate Up
create table device(
    device_eui bytea primary key,
    protocol_type varchar(50) not null default 'digital"',
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

-- +migrate Down
drop table device;