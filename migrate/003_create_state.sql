-- +migrate Up
create table device_state(
    device_eui bytea primary key,
    last_seen_at timestamp with time zone not null,
    location point,
    detail json
);

-- +migrate Down
drop table device_state;