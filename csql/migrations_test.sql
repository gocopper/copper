-- +migrate Up
create table people (name text);
insert into people (name) values ('test');

-- +migrate Down
drop table people;
