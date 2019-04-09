-- +migrate Up
create table device_track(
    device_eui bytea not null ,
    created_at timestamp with time zone not null,
    location point not null,
    altitude integer not null
);

create index idx_device_track on device_track(device_eui,created_at);
-- +migrate Down
drop index idx_device_track;
drop table device_track;