create table if not exists url_mappings
(
    id                 bigserial primary key,
    shortened_url      varchar(255)  unique not null,
    url                varchar(2048) not null,
    created_at         timestamp default timezone('UTC', now()),
    inactive_expire_at timestamp,
    must_expire_at     timestamp
);

create index if not exists url_mappings_shortened_url
    on url_mappings (shortened_url);