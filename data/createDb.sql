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
    teleId integer not null,
    year integer not null,
    month integer not null,
    sign integer not null,
    degree integer not null,
    FOREIGN KEY(teleId) REFERENCES user(teleId)
);

create table if not exists userVaccine(
    id integer primary key,
    created datetime,
    teleId integer not null,
    year integer not null,
    month integer not null,
    kind integer not null,
    effect integer not null,
    FOREIGN KEY(teleId) REFERENCES user(teleId)
);

create table if not exists userCountry(
    id integer primary key,
    rus string
);
insert into userCountry (rus) values ("Россия"), ("Украина"), ("Беларусь"), ("Казахстан");

create table if not exists userEducation(
    id integer primary key,
    rus string
);
insert into userEducation (rus) values ("Среднее"), ("Колледж"), ("Университет");

create table if not exists userVaccineOpinion(
    id integer primary key,
    rus string
);
insert into userVaccineOpinion (rus) values ("Помогают"), ("Бесполезны"), ("Опасны");

create table if not exists userOriginOpinion(
    id integer primary key,
    rus string
);
insert into userOriginOpinion (rus) values ("Природа"), ("Люди");

create table if not exists illnessSign(
    id integer primary key,
    rus string
);
insert into illnessSign (rus) values ("Есть медицинская справка"), ("Есть тест с наличием антител"), ("По характерным симптомам");

create table if not exists illnessDegree(
    id integer primary key,
    rus string
);
insert into illnessDegree (rus) values
    ("Лежал(а) под ИВЛ"), ("Лежал(а) в больнице"),
    ("Лежал(а) дома, тяжело"), ("Лежал(а) дома, средне"),
    ("Перенес(ла) на ногах"), ("Перенес(ла) без симптомов");

create table if not exists vaccineKind(
    id integer primary key,
    rus string
);
insert into vaccineKind (rus) values ("Спутник-V (два укола)"), ("Спутник-Лайт (один укол)"), ("ЭпиВакКорона"), ("КовиВак");

create table if not exists vaccineEffect(
    id integer primary key,
    rus string
);
insert into vaccineEffect (rus) values
    ("Сильные: температура, головная боль и т.п."),
    ("Средние: боль в руке, аллергия и т.п."),
    ("Слабые или никаких проявлений");

create table if not exists year(
    id integer primary key,
    rus string
);
insert into year (rus) values
    ("2020"), ("2021"), ("2022");

create table if not exists month(
    id integer primary key,
    rus string
);
insert into month (rus) values
    ("Январь"), ("Февраль"), ("Март"), ("Апрель"), ("Май"), ("Июнь"),
    ("Июль"), ("Август"), ("Сентябрь"), ("Октябрь"), ("Ноябрь"), ("Декабрь");
