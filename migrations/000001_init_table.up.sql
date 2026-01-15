create table if not exists users
(
    id         bigserial primary key,
    email      varchar(255) unique                      not null,
    password   varchar(255)                             not null,
    created_at timestamp default timezone('UTC', now()) not null,
    is_active  boolean   default true                   not null
);

create table if not exists urls
(
    id         bigserial primary key,
    user_id    bigint                                   null,
    title      varchar(255)                             not null,
    short_url  varchar(255)                             not null,
    full_url   varchar(2048)                            not null,
    created_at timestamp default timezone('UTC', now()) not null,
    is_deleted boolean   default false                  not null,
    constraint fk_user_id foreign key (user_id) references users (id)
        on delete set null on update cascade
);

create index if not exists idx_urls_user_id_created_at_active
    on urls (user_id, created_at desc)
    where is_deleted = false;

create unique index if not exists urls_short_url_active_unique
    on urls (short_url)
    where is_deleted = false;

create table if not exists redirects
(
    id         bigserial primary key,
    url_id     bigint                                   null,
    device     varchar(255)                             not null,
    country    varchar(255)                             not null,
    city       varchar(255)                             not null,
    created_at timestamp default timezone('UTC', now()) not null,
    constraint fk_url_id foreign key (url_id) references urls (id)
        on delete set null on update cascade
);

create index if not exists redirects_url_id
    on redirects (url_id);