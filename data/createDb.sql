create table if not exists user(
    teleId integer primary key,
    created datetime,
    modified datetime,
    name varchar(255),
    country integer not null,
    birth integer not null,
    gender integer not null,
    education integer not null,
    vaccineOpinion integer not null,
    originOpinion integer not null,
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

create table if not exists userCountry(
    id integer primary key,
    rus string,
    FOREIGN KEY(id) REFERENCES user(country)
);
insert into userCountry (rus) values ("Россия"), ("Украина"), ("Беларусь"), ("Казахстан");

create table if not exists userEducation(
    id integer primary key,
    rus string,
    FOREIGN KEY(id) REFERENCES user(education)
);
insert into userEducation (rus) values ("Среднее"), ("Колледж"), ("Университет");

create table if not exists userVaccineOpinion(
    id integer primary key,
    rus string,
    FOREIGN KEY(id) REFERENCES user(vaccineOpinion)
);
insert into userVaccineOpinion (rus) values ("Помогают"), ("Бесполезны"), ("Опасны");

create table if not exists userOriginOpinion(
    id integer primary key,
    rus string,
    FOREIGN KEY(id) REFERENCES user(originOpinion)
);
insert into userOriginOpinion (rus) values ("Природа"), ("Люди");
