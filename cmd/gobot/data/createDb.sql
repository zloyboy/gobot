create table if not exists user(
    teleId integer primary key,
    created datetime,
    modified datetime,
    country integer not null,
    birth integer not null,
    gender integer not null,
    education integer not null,
    vaccineOpinion integer not null,
    originOpinion integer not null,
    countIll integer not null,
    countVac integer not null
);

create table if not exists chat(
    id integer primary key
);

create table if not exists userAgeGroup(
    id integer primary key,
    created datetime,
    teleId integer not null,
    have_ill integer not null,
    have_vac integer not null,
    age_group integer not null,
    gender integer,
    FOREIGN KEY(teleId) REFERENCES user(teleId)
);

create table if not exists userIllness(
    id integer primary key,
    created datetime,
    teleId integer not null,
    year integer not null,
    month integer not null,
    sign integer not null,
    degree integer not null,
    age integer not null,
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
    age integer not null,
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
insert into userVaccineOpinion (rus) values ("Помогают"), ("Бесполезны"), ("Опасны"), ("Не знаю");

create table if not exists userOriginOpinion(
    id integer primary key,
    rus string
);
insert into userOriginOpinion (rus) values ("Природа"), ("Люди"), ("Не знаю");

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
    ("Критически: лежал(а) под ИВЛ"),
    ("Тяжело: лежал(а) в больнице"),
    ("Болел(а) дома: боль/температура"),
    ("Болел(а) дома: недомогание"),
    ("Перенес(ла) на ногах"),
    ("Перенес(ла) без симптомов");

create table if not exists vaccineKind(
    id integer primary key,
    rus string
);
insert into vaccineKind (rus) values ("Спутник-V (два укола)"), ("Спутник-Лайт"), ("ЭпиВакКорона"), ("КовиВак");

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
