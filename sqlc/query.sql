-- name: CreateUser :exec
insert into users (email, password)
values ($1, $2);

-- name: GetUserByEmail :one
select *
from users
where email = $1
  and is_active = true;

-- name: GetUser :one
select *
from users
where id = $1
  and is_active = true;

-- name: CreateUrl :one
insert into urls (user_id, title, short_url, full_url, qr)
values ($1, $2, $3, $4, $5)
returning *;

-- name: DeleteUrl :exec
update urls
set is_deleted = true
where id = $1;

-- name: GetUrl :one
select *
from urls
where id = $1
  and is_deleted = false;

-- name: GetUrlsByUserId :many
select *
from urls
where user_id = $1
  and is_deleted = false;

-- name: GetUserTotalClicks :one
select count(*)
from redirects r
where r.url_id in (select id
                   from urls
                   where user_id = $1
                     and is_deleted = false);

-- name: GetUrlByShortUrl :one
select *
from urls
where short_url = $1
  and is_deleted = false;

-- name: CreateRedirect :exec
insert into redirects (url_id, device, country, city)
values ($1, $2, $3, $4);

-- name: GetRedirectsByUrlId :many
select *
from redirects
where url_id = $1;