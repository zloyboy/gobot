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