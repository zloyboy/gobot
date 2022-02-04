create table if not exists user(
    teleId integer primary key,
    created datetime,
    modified datetime,
    name varchar(255),
    country varchar(255),
    birth integer not null,
    gender integer not null,
    education varchar(255),
    vaccine varchar(255),
    origin varchar(255),
    countIll integer not null,
    countVac integer not null
);