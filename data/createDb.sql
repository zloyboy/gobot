create table if not exists user(
    teleId integer primary key,
    created datetime,
    modified datetime,
    name varchar(255),
    country integer not null,
    birth integer not null,
    gender integer not null,
    education varchar(255),
    vaccine varchar(255),
    origin varchar(255),
    countIll integer not null,
    countVac integer not null
);

create index if not exists name_index on user (name);

create table if not exists userIllness(
    id integer primary key,
    created datetime,
    teleId integer,
    year integer,
    month integer,
    sign string,
    degree string,
    FOREIGN KEY(teleId) REFERENCES user(teleId)
);

create table if not exists userVaccine(
    id integer primary key,
    created datetime,
    teleId integer,
    year integer,
    month integer,
    kind string,
    effect string,
    FOREIGN KEY(teleId) REFERENCES user(teleId)
);

create table if not exists country(
    id integer primary key,
    rus string,
    FOREIGN KEY(id) REFERENCES user(country)
);

insert into country (rus) values ("Россия"), ("Украина"), ("Беларусь"), ("Казахстан");
