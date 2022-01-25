create table if not exists user(
    id integer primary key,
    created datetime,
    name varchar(255),
    age integer not null,
    res integer not null
);