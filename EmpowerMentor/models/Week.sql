CREATE TABLE week(
    id bigserial not null primary key,
    user_id bigint not null,
    monday int not null,
    tuesday int not null,
    wednesday int not null,
    thursday int not null,
    friday int not null,
    saturday int not null,
    sunday int not null,


    -- Foreign keys
    foreign key (user_id) references users(id)
);

CREATE TABLE meals (
    id bigserial not null primary key,
    image_type text,
    title varchar(255),
    ready_in_minutes int,
    servings int,
    source_url text,
    nutrient_id bigint not null,
    foreign key (nutrient_id) references nutrients(id)
);

create table nutrients(
    id bigserial not null primary key,
    calories double precision,
    protein double precision,
    fat double precision,
    carbohydrates double precision
);

create table monday(
    id bigserial not null primary key,
    user_id bigint not null,

    meal_ids bigint[] not null,
    nutrient_id bigint not null,

    foreign key (meal_ids) references meals(id),
    foreign key (nutrient_id) references nutrients(id)
);

create table tuesday(
    id bigserial not null primary key,
    user_id bigint not null,

    meal_ids bigint[] not null,
    nutrient_id bigint not null,

    foreign key (meal_ids) references meals(id),
    foreign key (nutrient_id) references nutrients(id)
);

create table wednesday(
    id bigserial not null primary key,
    user_id bigint not null,

    meal_ids bigint[] not null,
    nutrient_id bigint not null,

    foreign key (meal_ids) references meals(id),
    foreign key (nutrient_id) references nutrients(id)
);

create table thursday(
    id bigserial not null primary key,
    user_id bigint not null,

    meal_ids bigint[] not null,
    nutrient_id bigint not null,

    foreign key (meal_ids) references meals(id),
    foreign key (nutrient_id) references nutrients(id)
);

create table friday(
    id bigserial not null primary key,
    user_id bigint not null,

    meal_ids bigint[] not null,
    nutrient_id bigint not null,

    foreign key (meal_ids) references meals(id),
    foreign key (nutrient_id) references nutrients(id)
);

create table saturday(
    id bigserial not null primary key,
    user_id bigint not null,

    meal_ids bigint[] not null,
    nutrient_id bigint not null,

    foreign key (meal_ids) references meals(id),
    foreign key (nutrient_id) references nutrients(id)
);

create table sunday(
    id bigserial not null primary key,
    user_id bigint not null,

    meal_ids bigint[] not null,
    nutrient_id bigint not null,

    foreign key (meal_ids) references meals(id),
    foreign key (nutrient_id) references nutrients(id)
);

