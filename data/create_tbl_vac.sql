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